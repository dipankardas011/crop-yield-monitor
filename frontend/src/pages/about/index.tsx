import { Inter } from "next/font/google";


import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"



import { AlertCircle } from "lucide-react"

import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/components/ui/alert"

import { header } from "@/components/header";
import { footer } from "@/components/footer";
import { getTitle } from "@/components/page-title";


const inter = Inter({ subsets: ["latin"] });

export function AlertMessage(variant?: "default" | "destructive" | null | undefined, message?: string) {
  var title: string = ""
  if (variant === "destructive") {
    title = "Error"
  } else {
    title = "Heads up!"
  }

  return (
    <Alert variant={variant}>
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>
        {message}
      </AlertDescription>
    </Alert>
  )
}

export default function Home() {
  return (
    <main
      className={`flex min-h-screen flex-col items-center justify-between p-10 ${inter.className}`}
    >
      {header()}
      {getTitle("About Us")}
      <Accordion className="w-[400px]" type="single" collapsible>
        <AccordionItem value="project">
          <AccordionTrigger>What is this Project?</AccordionTrigger>
          <AccordionContent>
            Its a Major project POC
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="authors">
          <AccordionTrigger>Who all are the author(s)?</AccordionTrigger>
          <AccordionContent>
            Its a Major project POC
            <ul>
              <li>Dipankar Das</li>
              <li>Dipankar Das</li>
              <li>Dipankar Das</li>
              <li>Dipankar Das</li>
            </ul>
          </AccordionContent>
        </AccordionItem>
      </Accordion>

      {footer()}

    </main>
  );
}
