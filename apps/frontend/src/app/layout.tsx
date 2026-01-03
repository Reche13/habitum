import type { Metadata } from "next";
import Script from "next/script";
import { Poppins } from "next/font/google";
import "./globals.css";
import { QueryProvider } from "@/providers/query-provider";

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "Habitum",
  description: "Track your habits",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const googleClientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
  
  return (
    <html lang="en">
      <body className={`${poppins.className} antialiased`}>
        {googleClientId && (
          <Script
            src="https://accounts.google.com/gsi/client"
            strategy="afterInteractive"
          />
        )}
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
