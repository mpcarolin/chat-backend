import ChatBot from "react-chatbotify";
import type { Settings, Flow, Params } from "react-chatbotify";
import "./Chat.css";

type Path = "start" | "loop" | "end" | "error";

const STOP_COMMANDS = [
    "quit",
    "q",
    "end",
    "stop",
    "goodbye",
    "bye",
    "finish",
    "finished",
    "i'm finished",
    "i'm done",
    "done",
    "all done",
    "complete",
];


const chat = async (params: Params): Promise<Path> => {
    try {
        if (STOP_COMMANDS.includes(params.userInput.toLowerCase())) {
            return "end"
        }
        const res = await fetch("/api/chat", {
            method: "POST",
            body: JSON.stringify({ messages: [{ role: "user", content: params.userInput }] }),
            headers: {
                "Content-Type": "application/json"
            }
        });
        const json = await res.json();
        await params.injectMessage(json.response);
        return "loop"
    } catch (error) {
        console.error(error);
        return "error"
    }
}

export const Chat = ({ styles }: Parameters<typeof ChatBot>[0]) => {
    const blue = "#17262e"
    const orange = "#f54b06"
    const settings: Settings = {
        general: {
            showFooter: false,
            showHeader: true,
            primaryColor: orange,
            secondaryColor: blue,
            fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif",
        },
        notification: {
            disabled: true
        },
        header: {
            showAvatar: false,
            title: "Chat"
        },
        chatHistory: {
            storageKey: "chat-react-base-history"
        }
    }

    const flow: Flow = {
        start: {
            message: "Hi there, I'm here to help answer your questions about Star Wars. Go ahead and ask me anything!",
            path: chat
        },
        loop: {
            path: chat
        },
        error: {
            message: "I'm sorry, I'm having trouble answering that. Can you please try again?",
            path: "loop",
        },
        end: {
            message: "Okay! If you need anything else, just ask.",
            chatDisabled: false, // could make this true 
            path: chat,
        }
    }

    return (
        <div className="chatContainer">
            <ChatBot
                flow={flow}
                settings={settings}
                styles={styles}
            />
        </div>
    )
}
