import { ReactNode } from "react";

type ResponsiveSplitViewProps = {
  show?: boolean;
  desktop: ReactNode;
  mobile: ReactNode;
  desktopClassName?: string;
  mobileClassName?: string;
};

export function ResponsiveSplitView({
  show = true,
  desktop,
  mobile,
  desktopClassName = "hidden md:block",
  mobileClassName = "md:hidden",
}: ResponsiveSplitViewProps) {
  if (!show) {
    return null;
  }

  return (
    <>
      <div className={desktopClassName}>{desktop}</div>
      <div className={mobileClassName}>{mobile}</div>
    </>
  );
}
