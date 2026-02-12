// Highlight active nav link on scroll
const sections = document.querySelectorAll("section[id]");
const navLinks = document.querySelectorAll("nav ul a");

window.addEventListener("scroll", () => {
  const scrollY = window.scrollY + 120;

  for (const section of sections) {
    const top = section.offsetTop;
    const height = section.offsetHeight;
    const id = section.getAttribute("id");

    if (scrollY >= top && scrollY < top + height) {
      navLinks.forEach((link) => {
        link.classList.toggle(
          "active",
          link.getAttribute("href") === `#${id}`
        );
      });
    }
  }
});

// Chat form
const chatLog = document.getElementById("chat-log");

function formatReply(text) {
  // Sanitize HTML to prevent XSS
  text = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  // Bold: **text**
  text = text.replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>");
  // Italic: *text*
  text = text.replace(/\*(.+?)\*/g, "<em>$1</em>");
  // Convert list items and wrap consecutive ones in <ul>
  const lines = text.split("\n");
  const result = [];
  let inList = false;
  for (const line of lines) {
    const listMatch = line.match(/^[\-\*] (.+)/);
    if (listMatch) {
      if (!inList) { result.push("<ul>"); inList = true; }
      result.push("<li>" + listMatch[1] + "</li>");
    } else {
      if (inList) { result.push("</ul>"); inList = false; }
      result.push(line);
    }
  }
  if (inList) result.push("</ul>");
  return result.join("<br>");
}

function appendMsg(text, type) {
  const el = document.createElement("div");
  el.className = `chat-msg ${type}`;
  if (type === "sent") {
    el.textContent = "> " + text;
  } else {
    el.innerHTML = "&lt; " + formatReply(text);
  }
  chatLog.appendChild(el);
  chatLog.scrollTop = chatLog.scrollHeight;
}

document.getElementById("chat-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const input = document.getElementById("chat-input");
  const message = input.value.trim();
  if (!message) return;

  appendMsg(message, "sent");
  input.value = "";

  try {
    const res = await fetch("/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message }),
    });
    if (!res.ok) {
      console.error("chat error:", res.status);
      return;
    }
    const data = await res.json();
    if (data.reply) appendMsg(data.reply, "received");
  } catch (err) {
    console.error("chat fetch failed:", err);
  }
});
