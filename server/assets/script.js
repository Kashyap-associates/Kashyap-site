const data = {
  m_nav: window.matchMedia("(width < 700px)").matches,
  sdbr: false,
  chat: [
    {
      role: "system",
      content: "Hi There 👋! How can I help you today?",
    },
  ],

  updateNav(state) {
    this.m_nav = state;
  },

  talk(msg) {
    try {
      // pushing user data to chat
      this.chat.push({
        role: "user",
        content: msg,
      });

      // creating a copy of chat
      const body_data = [...this.chat];

      // chat animation (sortof...)
      setTimeout(() => {
        this.chat.push({
          role: "system",
          content: "...",
        });
      }, 300);

      setTimeout(async () => {
        // fetching data
        let data = null;
        const res = await fetch("/chat", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(body_data),
        });

        // parsing the response
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        } else {
          data = await res.json();
        }

        // replasing the content with that data
        if (data !== undefined || data !== null) {
          this.chat[this.chat.length - 1].content = data.msg;
        }
      }, 700);
    } catch (err) {
      console.error(err);
    }
  },
};

document.addEventListener("DOMContentLoaded", () => {
  const media = window.matchMedia("(max-width: 700px)");
  const mobile = document.getElementById("mobile");
  const desktop = document.getElementById("desktop");
  function handleResize(e) {
    if (mobile && desktop) {
      if (e.matches) {
        mobile.classList.remove("hide");
        desktop.classList.add("hide");
        data.updateNav(true);
      } else {
        mobile.classList.add("hide");
        desktop.classList.remove("hide");
        data.updateNav(false);
      }
    }
  }
  media.addEventListener("change", handleResize);
  handleResize(media);
});

function link_to(link) {
  if (link) {
    try {
      window.location.href = link;
    } catch (err) {
      console.error(err);
    }
  } else {
    console.error("Link is empty!");
  }
}
