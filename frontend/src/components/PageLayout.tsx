import { PropsWithChildren } from "react";

function PageLayout({ children }: PropsWithChildren) {
  return (
    <div className="p-4">
      <main>{children}</main>
    </div>
  );
}

export default PageLayout;
