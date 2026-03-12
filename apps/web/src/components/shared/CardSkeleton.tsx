import React from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";

interface CardSkeletonProps {
    cards?: number;
}

export function CardSkeleton({ cards = 4 }: CardSkeletonProps) {
    return (
        <div className="space-y-3">
            {Array.from({ length: cards }).map((_, i) => (
                <Card key={i}>
                    <CardContent className="py-4">
                        <Skeleton className="h-5 w-3/4 mb-2" />
                        <Skeleton className="h-4 w-1/2 mb-4" />
                        <Skeleton className="h-4 w-full mt-4" />
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}
