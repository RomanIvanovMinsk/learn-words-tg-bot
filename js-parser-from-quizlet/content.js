//parse words and answers
var node = document.querySelectorAll(".SetPageTerm-sideContent a.SetPageTerm-wordText span.TermText");
var description = document.querySelectorAll(".SetPageTerm-sideContent .SetPageTerm-definitionText span.TermText");

var wordsList = [];
if (node.keys.length == description.keys.length)
{
    
    node.forEach((item, key) => {
        wordsList.push({ word: item.textContent, value: description[key].textContent})
    })
    
    wordsList.forEach( (item) => {
        console.log(item)
    })
}

 var manifestData = chrome.runtime.getManifest()
 serviceUrl = manifestData.permissions[0];
 console.log(serviceUrl)
 // send them to the point of the bot that create list of words for the user
$.ajax({
    type: "POST",
    url: serviceUrl,
    data: data,
    success: function(msg){
      console.log("Success")
   },
   error: function(){
     console.log("error")
   }
 });