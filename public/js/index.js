fetch("/config/details.json")
  .then((res) => res.json())
  .then((data) => {
    document.querySelector("link[rel='icon']").href = data.icon;
    document.querySelector("nav .left #logo").src = data.icon;
    document.querySelector("#title").innerHTML = data.title;
    document.title = data.title;
  })
  .catch((err) =>
    console.log(`Error occured while fetching detils.json: ${err}`),
  );

fetch("/config/home.json")
  .then((res) => res.json())
  .then((data) => {
    document.querySelector("#home #tagline").innerHTML = data.tag_line;
    document.querySelector("#home #subtag").innerHTML = data.sub_tag;
  })
  .catch((err) =>
    console.log(`Error occured while fetching home.json: ${err}`),
  );

fetch("/config/about.json")
  .then((res) => res.json())
  .then((data) => {
    document.querySelector("#about .logo").src = data.about_icon;
    document.querySelector("#about #text").innerHTML = data.about_content;
  })
  .catch((err) =>
    console.log(`Error occured while fetching home.json: ${err}`),
  );

fetch("/config/services.json")
  .then((res) => res.json())
  .then((data) => {
    document.querySelector("#services .card .bottom p").innerHTML =
      data[0].title;
  })
  .catch((err) =>
    console.log(`Error occured while fetching home.json: ${err}`),
  );
