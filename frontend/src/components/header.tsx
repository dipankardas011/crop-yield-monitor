import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar"
import Link from "next/link"
import { ModeToggle } from "./theme-toogle"
import { Button } from "./ui/button"
import { Mail } from "lucide-react"


function getUserName(): string | undefined {
  return "dipankar"
}

function avatarDemo() {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Avatar>
          <AvatarImage src="https://github.com/dipankardas011.png" alt="@dipankardas011" />
          <AvatarFallback>DD</AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>{getUserName()}</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem>Logout</DropdownMenuItem>
        <DropdownMenuItem>Billing</DropdownMenuItem>
        <DropdownMenuItem>Team</DropdownMenuItem>
        <DropdownMenuItem>Subscription</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export function header() {

  return (<>
    <div className="border-b pb-10 mb-10">
      <div className="flex h-16 items-center px-4">
        <nav className="flex lg:space-x-6 mx-6">

          <div className="ml-auto flex items-center space-x-4 pl-20">
            <h2 className="text-2xl font-bold tracking-tight">Crop Yield Monitor</h2>
          </div>
          <div className="ml-auto flex items-center space-x-4">
            <Button variant="outline" className="dark:bg-transparent text-5l font-bold" asChild>
              <Link href="/">Home</Link>
            </Button>

            <Button variant="outline" className="dark:bg-transparent text-5l font-bold" asChild>
              <Link href="/about">About Us</Link>
            </Button>
          </div>

          <div className="ml-auto flex items-center space-x-4 pl-10">
            {avatarDemo()}
            {ModeToggle()}
          </div>
        </nav>
      </div>
    </div>
  </>
  );
}