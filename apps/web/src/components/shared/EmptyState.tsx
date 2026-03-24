import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { LucideIcon } from "lucide-react";

interface EmptyStateProps {
    icon: LucideIcon;
    title: string;
    description?: string;
    action?: React.ReactNode;
}

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
    return (
        <Card className="border-dashed shadow-sm">
            <CardContent className="flex flex-col items-center justify-center py-16 text-center">
                <div className="rounded-full bg-muted p-4 mb-4">
                    <Icon className="h-8 w-8 text-muted-foreground" />
                </div>
                <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">{title}</h3>
                {description && (
                    <p className="mt-2 text-sm text-muted-foreground max-w-sm">
                        {description}
                    </p>
                )}
                {action && <div className="mt-6">{action}</div>}
            </CardContent>
        </Card>
    );
}
