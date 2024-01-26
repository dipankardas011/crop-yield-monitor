import { Inter } from "next/font/google";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button";
import { header } from "@/components/header";
import { footer } from "@/components/footer";
import { AlertMessage } from "@/components/pop-message";
import { getTitle } from "@/components/page-title";

const inter = Inter({ subsets: ["latin"] });

export function TextareaPrediction() {
  return (
    <div className="grid w-full gap-2">
      <Textarea placeholder="Predictions will come here" disabled />
      <Button>Generate</Button>
    </div>
  );
}

export function InputFile() {
  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="picture">Picture</Label>
      <Input id="picture" type="file" />
      {AlertMessage("destructive", "dummy error")}
      {AlertMessage("default", "dummy info")}
    </div>
  )
}


export default function Home() {
  return (
    <main
      className={`flex min-h-screen flex-col items-center justify-between p-10 ${inter.className}`}
    >
      {header()}
      {getTitle("Dashboard")}

      <Tabs defaultValue="soil_image" className="w-[400px]">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="soil_image">Soil Image</TabsTrigger>
          <TabsTrigger value="predictions">Generate Prediction</TabsTrigger>
        </TabsList>
        <TabsContent value="soil_image">{InputFile()}</TabsContent>
        <TabsContent value="predictions">{TextareaPrediction()}</TabsContent>
      </Tabs>

      {footer()}
    </main>
  );
}