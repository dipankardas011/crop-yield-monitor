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
    },
  })

  // 2. Define a submit handler.
  function onSubmit(values: z.infer<typeof formSchema>) {

    const container = document.getElementById('response-login');

    const root = createRoot(container!);

    async function login() {

      try {
        const response = await fetch('http://localhost:9090/account/signin', {
          method: 'POST', // Change to POST if necessary
          credentials: 'include',
          body: JSON.stringify({ username: values.username, password: values.password }), // Uncomment for POST
        });

        const data = await response.json();

        // Check for successful login (adjust based on your backend response structure)
        if (response.ok) {
          // Redirect to the dashboard
          root.render(AlertMessage("default", "Login successfully\n" + JSON.stringify(data, null, 2)));
          window.location.href = '/account/dashboard/';
        } else {
          // Handle login failure (display an error message, etc.)
          root.render(AlertMessage("destructive", "Login failed\n" + String(data.error)));
          console.error('Login failed:', data.error);
        }
      } catch (error) {
        console.error('Error during login:', error);
        root.render(AlertMessage("destructive", "Login failed\n" + String(error)));
        console.error('Login failed:', error);
      }
    }
    login()
  }

  return (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 p-10">
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
      <div id="response-login"></div>
    </>
  );
}