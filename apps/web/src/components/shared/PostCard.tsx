import React from "react";
import { Post, PostComment } from "@/types";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { POST_SCOPE_LABELS, POST_TYPE_META } from "@/lib/post-config";
import { Heart, MessageCircle, SendHorizontal } from "lucide-react";
import { teacherApi } from "@/lib/api/teacher.api";
import { parentApi } from "@/lib/api/parent.api";

interface PostCardProps {
  post: Post;
  authorLabel?: string;
  enableInteractions?: boolean;
  enableShare?: boolean;
  audience?: "teacher" | "parent";
  onPostPatched?: (postId: string, patch: Partial<Post>) => void;
}

function getScopeDisplay(post: Post): string {
  return POST_SCOPE_LABELS[post.scope_type] || post.scope_type;
}

export function PostCard({
  post,
  authorLabel = "Giáo viên",
  enableInteractions = true,
  enableShare = true,
  audience = "teacher",
  onPostPatched,
}: PostCardProps) {
  const [liked, setLiked] = React.useState(post.liked_by_me);
  const [likeCount, setLikeCount] = React.useState(post.like_count || 0);
  const [shareCount, setShareCount] = React.useState(post.share_count || 0);
  const [commentCount, setCommentCount] = React.useState(post.comment_count || 0);
  const [showComments, setShowComments] = React.useState(false);
  const [commentDraft, setCommentDraft] = React.useState("");
  const [comments, setComments] = React.useState<PostComment[]>([]);
  const [commentsLoaded, setCommentsLoaded] = React.useState(false);
  const [loadingComments, setLoadingComments] = React.useState(false);
  const [submittingComment, setSubmittingComment] = React.useState(false);
  const [processingLike, setProcessingLike] = React.useState(false);
  const [processingShare, setProcessingShare] = React.useState(false);
  const [interactionError, setInteractionError] = React.useState("");

  const postApi = audience === "parent" ? parentApi : teacherApi;

  React.useEffect(() => {
    setLiked(post.liked_by_me);
    setLikeCount(post.like_count || 0);
    setShareCount(post.share_count || 0);
    setCommentCount(post.comment_count || 0);
  }, [post.liked_by_me, post.like_count, post.share_count, post.comment_count]);

  const postMeta = POST_TYPE_META[post.type] || {
    label: post.type,
    badgeVariant: "secondary" as const,
  };

  const patchPost = React.useCallback(
    (patch: Partial<Post>) => {
      onPostPatched?.(post.post_id, patch);
    },
    [onPostPatched, post.post_id],
  );

  const handleLikeToggle = async () => {
    try {
      setProcessingLike(true);
      setInteractionError("");
      const payload = await postApi.togglePostLike(post.post_id);
      setLiked(payload.liked_by_me);
      setLikeCount(payload.like_count);
      patchPost({ liked_by_me: payload.liked_by_me, like_count: payload.like_count });
    } catch {
      setInteractionError("Không thể cập nhật lượt thích. Vui lòng thử lại.");
    } finally {
      setProcessingLike(false);
    }
  };

  const loadComments = React.useCallback(async () => {
    try {
      setLoadingComments(true);
      setInteractionError("");
      const response = await postApi.getPostComments(post.post_id, { limit: 50, offset: 0 });
      setComments(response.data || []);
      setCommentsLoaded(true);
    } catch {
      setInteractionError("Không thể tải bình luận. Vui lòng thử lại.");
    } finally {
      setLoadingComments(false);
    }
  }, [post.post_id, postApi]);

  const handleShare = async () => {
    try {
      setProcessingShare(true);
      setInteractionError("");
      const payload = await postApi.sharePost(post.post_id);
      setShareCount(payload.share_count);
      patchPost({ share_count: payload.share_count });
    } catch {
      setInteractionError("Không thể chia sẻ bài viết. Vui lòng thử lại.");
    } finally {
      setProcessingShare(false);
    }
  };

  const handleCommentSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    const nextComment = commentDraft.trim();
    if (!nextComment) {
      return;
    }

    try {
      setSubmittingComment(true);
      setInteractionError("");
      const payload = await postApi.createPostComment(post.post_id, { content: nextComment });
      setComments((prev) => [payload.comment, ...prev]);
      setCommentCount(payload.comment_count);
      patchPost({ comment_count: payload.comment_count });
      setCommentDraft("");
    } catch {
      setInteractionError("Không thể gửi bình luận. Vui lòng thử lại.");
    } finally {
      setSubmittingComment(false);
    }
  };

  const handleToggleComments = async () => {
    const next = !showComments;
    setShowComments(next);

    if (next && !commentsLoaded) {
      await loadComments();
    }
  };

  return (
    <Card className="overflow-hidden">
      <CardContent className="space-y-3 p-4 sm:p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="flex min-w-0 items-center gap-3">
            <Avatar>
              <AvatarFallback>GV</AvatarFallback>
            </Avatar>
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold text-foreground">{authorLabel}</p>
              <p className="text-xs text-muted-foreground">
                {new Date(post.created_at).toLocaleString("vi-VN")} • {getScopeDisplay(post)}
              </p>
            </div>
          </div>
          <Badge variant={postMeta.badgeVariant}>{postMeta.label}</Badge>
        </div>

        <p className="whitespace-pre-line text-sm leading-6 text-foreground">{post.content}</p>

        {enableInteractions && (
          <>
            <div className="flex flex-wrap items-center gap-2 border-t pt-3">
              <Button
                type="button"
                size="sm"
                variant="ghost"
                className={liked ? "text-primary" : "text-muted-foreground"}
                onClick={handleLikeToggle}
                aria-pressed={liked}
                disabled={processingLike}
              >
                <Heart className="mr-1 h-4 w-4" />
                Thích ({likeCount})
              </Button>

              <Button
                type="button"
                size="sm"
                variant="ghost"
                className="text-muted-foreground"
                onClick={handleToggleComments}
                aria-expanded={showComments}
              >
                <MessageCircle className="mr-1 h-4 w-4" />
                Bình luận ({commentCount})
              </Button>

              {enableShare && (
                <Button
                  type="button"
                  size="sm"
                  variant="ghost"
                  className="text-muted-foreground"
                  onClick={handleShare}
                  disabled={processingShare}
                >
                  <SendHorizontal className="mr-1 h-4 w-4" />
                  Chia sẻ ({shareCount})
                </Button>
              )}
            </div>

            {interactionError && <p className="text-xs text-destructive">{interactionError}</p>}

            {showComments && (
              <div className="space-y-3 border-t pt-3">
                <form className="flex items-center gap-2" onSubmit={handleCommentSubmit}>
                  <Input
                    value={commentDraft}
                    onChange={(event) => setCommentDraft(event.target.value)}
                    placeholder="Viết bình luận..."
                    aria-label="Viết bình luận"
                  />
                  <Button type="submit" size="sm" disabled={!commentDraft.trim() || submittingComment}>
                    Gửi
                  </Button>
                </form>

                {loadingComments && <p className="text-xs text-muted-foreground">Đang tải bình luận...</p>}

                {!loadingComments && comments.length > 0 && (
                  <ul className="space-y-2">
                    {comments.map((comment) => (
                      <li key={comment.comment_id} className="rounded-md bg-muted/50 px-3 py-2 text-sm text-foreground">
                        <p className="text-xs font-medium text-muted-foreground">{comment.author_display}</p>
                        <p>{comment.content}</p>
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
