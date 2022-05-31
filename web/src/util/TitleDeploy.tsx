import { useEffect } from "react";
import { TitleBase } from "./shared";

type Properties = {
  title?: string,
  children?: JSX.Element
}

export function TitleDeploy({ title, children }: Properties) {
  useEffect(() => {
    if (!document) {
      return;
    }
    document.title = `${TitleBase} ${title || "Website"}`;
  }, [title]);
  
  return (
    <>
      {children}
    </>
  );
}