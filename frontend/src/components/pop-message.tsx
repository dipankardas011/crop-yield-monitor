import {
    Alert,
    AlertDescription,
    AlertTitle,
} from "@/components/ui/alert"
import { AlertCircle, AlertTriangle } from "lucide-react"

export function AlertMessage(variant?: "default" | "destructive" | null | undefined, message?: string) {
    if (variant === "destructive") {
        return (
            <Alert variant={variant}>
                <AlertTriangle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                    {message}
                </AlertDescription>
            </Alert>
        )
    }

    return (
        <Alert variant={variant}>
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Heads Up!</AlertTitle>
            <AlertDescription>
                {message}
            </AlertDescription>
        </Alert>
    )
}
