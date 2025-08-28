import React from "react";
import "./Message.css";

/**
 * Individual chat message bubble with positioning and loading state support
 */
export const Message = React.forwardRef<HTMLDivElement, {
    position: "left" | "right",
    content?: React.ReactNode,
    loading?: boolean
}>(({ position, content, loading }, ref) => {
    const messageClass = position === "left" ? "system-message" : "user-message"
    const loadingClass = loading ? "message-shimmer" : ""
    return (
        <div ref={ref} className={`message ${messageClass} ${loadingClass}`}>
            {!loading && <p>{content}</p>}
        </div>
    )
});

Message.displayName = "Message";