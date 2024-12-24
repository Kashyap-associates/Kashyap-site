fetch("config.json")
  .then((res) => res.json())
  .then((data) => {
    document.title = data.title;
    document.querySelector("span#title").textContent = data.title;
    document.querySelector("a#instagram").href = data.social_media.instagram;
    document.querySelector("a#linkedin").href = data.social_media.linkedin;
    document.querySelector("a#whatsapp").href = data.social_media.whatsapp;
    document.querySelector("a#gmail").href = data.social_media.gmail;
  })
  .catch((err) => console.error("Error loading content :", err));
