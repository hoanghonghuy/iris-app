import React from 'react';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { Loader2 } from "lucide-react";

export interface ConfirmAlertDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  description: React.ReactNode;
  loading?: boolean;
  cancelText?: string;
  confirmText?: string;
}

export function ConfirmAlertDialog({
  isOpen, onClose, onConfirm, title, description, loading, cancelText = "Cancel", confirmText = "Confirm"
}: ConfirmAlertDialogProps) {
  return (
    <AlertDialog open={isOpen} onOpenChange={(open: boolean) => !open && onClose()}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription className="text-sm">{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={(e: React.MouseEvent<HTMLButtonElement>) => { e.preventDefault(); onClose(); }} disabled={loading}>{cancelText}</AlertDialogCancel>
          <AlertDialogAction onClick={(e: React.MouseEvent<HTMLButtonElement>) => { e.preventDefault(); onConfirm(); }} disabled={loading} className="bg-destructive hover:bg-destructive/90 text-destructive-foreground">
            {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : confirmText}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
