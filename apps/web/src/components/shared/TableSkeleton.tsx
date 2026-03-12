import React from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";

interface TableSkeletonProps {
    columns?: number;
    rows?: number;
}

export function TableSkeleton({ columns = 4, rows = 5 }: TableSkeletonProps) {
    return (
        <Card>
            <CardContent className="p-0">
                <table className="w-full">
                    <thead>
                        <tr className="border-b text-left text-sm text-muted-foreground">
                            {Array.from({ length: columns }).map((_, i) => (
                                <th key={i} className="px-6 py-4 font-medium">
                                    <Skeleton className="h-4 w-24" />
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {Array.from({ length: rows }).map((_, rIndex) => (
                            <tr key={rIndex} className="border-b last:border-0 hover:bg-zinc-50">
                                {Array.from({ length: columns }).map((_, cIndex) => (
                                    <td key={cIndex} className="px-6 py-4">
                                        <Skeleton className="h-4 w-full max-w-[200px]" />
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </CardContent>
        </Card>
    );
}
