import { Inter } from "next/font/google";



import { header } from "@/components/header";
import { footer } from "@/components/footer";
import { Button } from "@/components/ui/button";
import Link from "next/link";



const inter = Inter({ subsets: ["latin"] });


export default function Home() {
  return (
    <main
      className={`flex min-h-screen flex-col items-center justify-between p-10 ${inter.className}`}
    >
      {header()}

      <div className="flex items-center justify-between px-10 py-10">
        <div className="text-5xl font-bold">Welcome to <span className="home-page">Crop Yield Monitor</span> !</div>
      </div>
      <Button asChild>
        <Link href="/account/dashboard">Dashboard</Link>
      </Button>

      {footer()}
    </main>
  );
}


