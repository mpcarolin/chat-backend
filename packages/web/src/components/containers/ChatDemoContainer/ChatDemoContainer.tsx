import { Chat } from "../../common/Chat/Chat"
import { systemMessage, useChat, userMessage } from "../../../hooks/useChat"
import { shouldEnableStreaming } from "../../../utils/chatConfig"
import "./ChatDemoContainer.css"

/**
 * Container that manages chat state and provides message handling for the chat demo
 */
export const ChatDemoContainer = () => {
    const { messages, sendMessage, loading } = useChat({
        initialMessages: [
            systemMessage("How can I help you?"),
        ],
    });

    const enableStreaming = shouldEnableStreaming();

    const handleMessageSubmit = (message: string) => {
        sendMessage(
            userMessage(message),
            { stream: enableStreaming }
        );
    };

    return (
        <div className="chat-demo-container">
            <Chat
                messages={messages}
                onMessageSubmit={handleMessageSubmit}
                loading={loading}
            />
        </div>
    )
}
