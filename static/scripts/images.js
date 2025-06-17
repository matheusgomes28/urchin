// MARK: Images Modal

// TODO : This image stuff is probably better suited in an
// TODO : image manager singleton class

/**
 * Holds the array to image information in
 * current page
 * [{name: "1238-1023123-wer123123.jpg", "excerpt": "some test about my image", src="http://..."}, ...]
 */
const pageImages = []

/**
 * index to current image in image array
 */
let currImage = 0;

/**
 * Shows the modal image popup
 * 
 * @param {string} index index of the image in current page
 */
function showImageModal(index) {

  // Update the highlighting
  for (let i = 0; i < pageImages.length; ++i) {
    document.getElementById(`modal-pagination-${i}`).classList.remove(`text-indigo-500`);
  }
  document.getElementById(`modal-pagination-${index}`).classList.add(`text-indigo-500`);

  currImage = index;
  const image = pageImages[currImage];

  document.getElementById("modal-title").innerHTML = image.name;
  document.getElementById("modal-excerpt").innerHTML = image.excerpt;
  document.getElementById("modal-image").src = image.filepath;
  document.getElementById("modal-text-date").innerHTML = image.date;
  document.getElementById("modal-text-name").innerHTML = image.location.name;
  modal.showModal();
}

function nextImage() {
  currImage = (currImage + 1) % pageImages.length;
  showImageModal(currImage);
}

function prevImage() {
  currImage = (currImage - 1 + pageImages.length) % pageImages.length;
  showImageModal(currImage);
}

function refreshImages(images) {
  if (!Array.isArray(images)) {
    console.error("input must be an array of valid images");
  }

  // TODO : Currently, this will input some dummy data in excerpt
  // TODO : So we need to change the golang Image repr

  const imageData = images.filter(im => {
    const hasProperties = im.hasOwnProperty("uuid")
      && (typeof im.uuid === "string")
      && im.hasOwnProperty("filename")
      && (typeof im.name === "string")
      && im.hasOwnProperty("excerpt")
      && (typeof im.excerpt === "string")
      && im.hasOwnProperty("date")
      && (typeof im.date === "string")
      && im.hasOwnProperty("location")
      && (typeof im.location === "object")
      && im.location.hasOwnProperty("latitude")
      && (typeof im.location.latitude === "Number")
      && im.location.hasOwnProperty("longitude")
      && (typeof im.location.longitude === "Number")
      && im.location.hasOwnProperty("name")
      && (typeof im.location.name === "string");

    return !hasProperties;
  });

  const sanitizedImages = imageData.map(im => {
    const newImage = {...im};
    newImage.excerpt = im.excerpt === "" ? im.name : im.excerpt;
    return newImage;
  });

  pageImages.push(...sanitizedImages);
}
