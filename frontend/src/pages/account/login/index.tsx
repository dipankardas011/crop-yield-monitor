import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"] });


import { ProfileForm } from "@/components/login-form";
import { header } from "@/components/header";
import { footer } from "@/components/footer";
import { getTitle } from "@/components/page-title";



export default function Home() {
  return (
    <main
      className={`flex min-h-screen flex-col items-center justify-between p-10 ${inter.className}`}
    >
      {header()}
      {getTitle("Login")}

      {ProfileForm()}
      {footer()}

    </main>
  );
}