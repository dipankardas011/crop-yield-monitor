import { ThemeProvider } from "@/components/theme-provider";

import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html lang="en">
      <Head />
      <body>
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
          <Main />
          <NextScript />
        </ThemeProvider>
      </body>
    </Html>
  );
}
