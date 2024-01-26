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
import { createRoot } from 'react-dom/client'

const inter = Inter({ subsets: ["latin"] });

function TextareaPrediction() {
  return (
    <div className="grid w-full gap-2">
      <Button onClick={Recommendation}>Generate</Button>
      <div id="predictions"></div>
    </div>
  );
}

function InputFile() {
  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="picture">Picture</Label>
      <Input id="picture" type="file" onChange={ImageFile} />
      <div id="image-file"></div>
    </div>
  )
}

function ImageFile(event: any) {
  const container = document.getElementById('image-file');

  const root = createRoot(container!);


  const file = event.target.files[0];
  let format = file.type;
  if (file.type === "image/png" || file.type === "image/jpeg") {
    console.log("Image format");
  } else {
    console.log("Image format: Unknown");
    alert("Invalid file");
    return;
  }

  const reader = new FileReader();

  reader.onloadend = async function () {
    const arrayBuffer = reader.result as ArrayBuffer;
    const byteArray = new Uint8Array(arrayBuffer);
    const byteValues = Array.from(byteArray);
    const upload = async () => {
      try {
        const response = await fetch('http://localhost:9090/image/upload', {
          method: 'POST',
          credentials: 'include',
          body: JSON.stringify({ raw_image_bytes: byteValues, image_format: String(format) })
        });

        const data = await response.json();

        if (response.ok) {
          root.render(AlertMessage("default", "Image Uploaded\n" + JSON.stringify(data, null, 2)));
        } else {
          root.render(AlertMessage("destructive", "Image Upload failed\n" + data.error));
          console.error('Image Upload failed:', data.error);
        }
      } catch (error) {
        root.render(AlertMessage("destructive", "Image Upload failed\n" + String(error)));
        console.error('Image Upload failed:', String(error));
      }
    }
    upload()
  }
  reader.readAsArrayBuffer(file);
}

function Recommendation() {
  const container = document.getElementById('predictions');

  const root = createRoot(container!);

  const recommend = async () => {

    try {
      const response = await fetch('http://localhost:9090/recommend/get', {
        method: 'GET',
        credentials: 'include',
      });

      const data = await response.json();

      if (response.ok) {
        console.log(data);
        root.render(AlertMessage("default", JSON.stringify(data.stdout, null, 2) + "; " + JSON.stringify(data.recommendations, null, 2)));
      } else {
        root.render(AlertMessage("destructive", "recommend failed\n" + String(data.error)));
        console.error('recommend failed:', String(data.error));
      }
    } catch (error) {
      root.render(AlertMessage("destructive", "recommend failed\n" + String(error)));
      console.error('recommend failed:', String(error));
    }
  }
  recommend()
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