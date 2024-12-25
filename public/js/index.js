fetch("/config/details.json")
  .then((res) => res.json())
  .then((data) => {
    document.querySelector("link[rel='icon']").href = data.icon;
    document.querySelector("nav .top #logo").src = data.icon;
    document.querySelector("#title").innerHTML = data.title;
    document.title = data.title;
  })
  .catch((err) => console.log(`Error occured : ${err}`));

fetch("/config/home.json");

fetch("/config/about.json");

fetch("/config/services.json");
