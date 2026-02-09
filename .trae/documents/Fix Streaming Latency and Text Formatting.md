I have identified the causes for the streaming delay and formatting issues:

1.  **Streaming Delay ("Terminal answers first"):**
    *   **Cause:** Nginx (the web server) is likely buffering the backend's response. Even though the backend is generating text, Nginx waits for a certain amount of data before sending it to the browser.
    *   **Fix:** I will add a special header `X-Accel-Buffering: no` in the Go backend to explicitly tell Nginx to disable buffering for these requests. I will also optimize the backend flushing logic to send chunks immediately rather than character-by-character (which can sometimes trigger buffering behavior).

2.  **Formatting "Piled Up" (No newlines):**
    *   **Cause:** When the Markdown parser encounters incomplete text (during streaming) or fails, it falls back to displaying raw text. Without the specific CSS property `white-space: pre-wrap`, HTML collapses all newlines and spaces into a single space, causing the text to "pile up".
    *   **Fix:** I will add `white-space: pre-wrap` to the chat message CSS. This ensures that even raw text preserves line breaks and spacing perfectly.

**Implementation Plan:**

1.  **Backend (`service/session/session.go`):**
    *   Add `w.Header().Set("X-Accel-Buffering", "no")`.
    *   Simplify the loop to flush data chunks directly instead of character-by-character.

2.  **Frontend (`vue-frontend/src/views/AIChat.vue`):**
    *   Add `white-space: pre-wrap` to `.message-content` CSS.
    *   Ensure the Markdown rendering logic is robust.

3.  **Rebuild:**
    *   Rebuild the containers to apply these changes.

This will ensure immediate streaming response and correct formatting at all times.