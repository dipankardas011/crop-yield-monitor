"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"

import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { AlertMessage } from "./pop-message"
import { createRoot } from 'react-dom/client'


const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  name: z.string().min(5, {message: "Name must be at least 5 characters."}),
  emailid: z.string().min(5, {message: "Email must be at least 5 characters"}),
  password: z.string().min(8, {
    message: "Password must be at least 8 characters.",
  }),
})


export function ProfileForm() {
  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
      name: "",
      emailid: "",
    },
  })

  function onSubmit(values: z.infer<typeof formSchema>) {

    const container = document.getElementById('response-signup');

    const root = createRoot(container!);

    async function signup() {

      try {
        const response = await fetch('http://localhost:9090/account/signup', {
          method: 'POST',
          body: JSON.stringify({
            name: values.name,
            username: values.username,
            password: values.password,
            email: values.emailid,
          }),
        });

        const data = await response.json();

        // Check for successful signup (adjust based on your backend response structure)
        if (response.ok) {
          // Redirect to the login page or perform other actions as needed
          root.render(AlertMessage("default", "Signup successfully\n" + JSON.stringify(data, null, 2)));
        } else {
          // Handle signup failure (display an error message, etc.)
          root.render(AlertMessage("destructive", "Signup failed\n" + data.error));
          console.error('Signup failed:', data.error);
        }
      } catch (error) {
        root.render(AlertMessage("destructive", "Signup failed\n" + String(error)));
        console.error('Signup failed:', error);
      }
    }
    signup()
  }

  return (
    <>
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5 p-10">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="full name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="admin" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="emailid"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="email id" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input placeholder="XXX" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
    <div id="response-signup"></div>
    </>
  )
}
