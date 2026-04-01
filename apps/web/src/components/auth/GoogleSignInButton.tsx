"use client";

import { useEffect, useState, useRef } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

declare global {
  interface Window {
    google?: {
      accounts?: {
        id?: {
          initialize: (config: {
            client_id: string;
            callback: (response: { credential: string }) => void;
            ux_mode?: "popup" | "redirect";
          }) => void;
          renderButton: (
            parent: HTMLElement,
            options: Record<string, string | number | boolean>
          ) => void;
          prompt: () => void;
        };
      };
    };
  }
}

type Props = {
  onSubmitGoogle: (payload: { idToken: string; password?: string }) => Promise<void>;
  /** error_code machine-readable từ backend (dùng để detect trạng thái, không phải text) */
  errorCode?: string;
  clearError: () => void;
  disabled?: boolean;
};

const GOOGLE_SCRIPT_SRC = "https://accounts.google.com/gsi/client";

export function GoogleSignInButton({ onSubmitGoogle, errorCode, clearError, disabled }: Props) {
  const [scriptReady, setScriptReady] = useState(false);
  const [pendingCredential, setPendingCredential] = useState("");
  const [linkPassword, setLinkPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const initialized = useRef(false);
  const oneTapPrompted = useRef(false);

  const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";
  const isUiBlocked = Boolean(disabled || isSubmitting);

  useEffect(() => {
    if (!clientId) {
      return;
    }

    const existing = document.querySelector(`script[src='${GOOGLE_SCRIPT_SRC}']`) as HTMLScriptElement | null;
    if (existing) {
      if (window.google?.accounts?.id) {
        setScriptReady(true);
      } else {
        existing.addEventListener("load", () => setScriptReady(true), { once: true });
      }
      return;
    }

    const script = document.createElement("script");
    script.src = GOOGLE_SCRIPT_SRC;
    script.async = true;
    script.defer = true;
    script.onload = () => setScriptReady(true);
    document.head.appendChild(script);
  }, [clientId]);

  useEffect(() => {
    if (!scriptReady || !clientId || !window.google?.accounts?.id) {
      return;
    }

    const buttonRoot = document.getElementById("google-signin-render-root");
    if (!buttonRoot) {
      return;
    }
    buttonRoot.innerHTML = "";
    const isDarkTheme = document.documentElement.classList.contains("dark");
    const buttonWidth = Math.min(360, Math.max(240, Math.floor(buttonRoot.clientWidth || 320)));

    if (!initialized.current) {
      window.google.accounts.id.initialize({
        client_id: clientId,
        ux_mode: "popup",
        callback: async (response) => {
          clearError();
          setPendingCredential(response.credential);
          setIsSubmitting(true);
          try {
            await onSubmitGoogle({ idToken: response.credential });
          } finally {
            setIsSubmitting(false);
          }
        },
      });
      initialized.current = true;
    }

    if (isUiBlocked) {
      buttonRoot.innerHTML = "";
      return;
    }

    window.google.accounts.id.renderButton(buttonRoot, {
      type: "standard",
      theme: isDarkTheme ? "filled_black" : "outline",
      size: "large",
      shape: "rectangular",
      text: "signin_with",
      width: buttonWidth,
    });
  }, [clearError, clientId, isUiBlocked, onSubmitGoogle, scriptReady]);

  useEffect(() => {
    if (!scriptReady || !clientId || !window.google?.accounts?.id) {
      return;
    }
    if (isUiBlocked || pendingCredential || oneTapPrompted.current) {
      return;
    }

    window.google.accounts.id.prompt();
    oneTapPrompted.current = true;
  }, [clientId, isUiBlocked, pendingCredential, scriptReady]);

  const showPasswordLinkStep = Boolean(
    pendingCredential && errorCode === "GOOGLE_LINK_PASSWORD_REQUIRED"
  );

  const handleLinkAndLogin = async () => {
    if (!pendingCredential || !linkPassword || isSubmitting) {
      return;
    }
    setIsSubmitting(true);
    clearError();
    try {
      await onSubmitGoogle({ idToken: pendingCredential, password: linkPassword });
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!clientId) {
    return (
      <p className="text-xs text-muted-foreground text-center">
        Google Sign-In chưa được cấu hình cho môi trường này.
      </p>
    );
  }

  return (
    <div className="space-y-3 w-full">
      <div
        className="w-full flex justify-center"
        id="google-signin-render-root"
        aria-busy={isUiBlocked}
      />

      {showPasswordLinkStep && (
        <div className="space-y-2 rounded-md border p-3">
          <p className="text-xs text-muted-foreground">
            Tài khoản chưa liên kết Google. Nhập mật khẩu hiện tại để xác nhận liên kết.
          </p>
          <Input
            type="password"
            value={linkPassword}
            onChange={(e) => setLinkPassword(e.target.value)}
            placeholder="Mật khẩu hiện tại"
            disabled={isSubmitting || disabled}
          />
          <Button
            type="button"
            variant="secondary"
            className="w-full"
            onClick={handleLinkAndLogin}
            disabled={isSubmitting || disabled || !linkPassword}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Đang liên kết...
              </>
            ) : (
              "Xác nhận liên kết và đăng nhập"
            )}
          </Button>
        </div>
      )}
    </div>
  );
}
