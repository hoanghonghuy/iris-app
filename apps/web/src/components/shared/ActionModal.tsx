import React from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";

export interface ActionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  description?: React.ReactNode;
  loading?: boolean;
  disabled?: boolean;
  confirmText?: string;
  cancelText?: string;
  children?: React.ReactNode;
}

export function ActionModal({
  isOpen, onClose, onConfirm, title, description, loading, disabled, confirmText = "Confirm", cancelText = "Cancel", children
}: ActionModalProps) {
  return (
    <Dialog open={isOpen} onOpenChange={(open: boolean) => !open && onClose()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && <DialogDescription className="text-sm">{description}</DialogDescription>}
        </DialogHeader>
        {children && (
          <div className="py-2">
            {children}
          </div>
        )}
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={loading}>{cancelText}</Button>
          <Button onClick={onConfirm} disabled={loading || disabled}>
            {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} {confirmText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
