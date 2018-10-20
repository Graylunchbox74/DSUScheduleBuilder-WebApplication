$(function() {
    const SEARCH_STRING = "/courses/search"
    function getCoursesFromSearchTerms() {
        fetch(SEARCH_STRING, {
            method: "POST",
            cache: "no-cache",
            headers: {
                "Content-Type": "application/json; charset=utf-8"
            },
            body: JSON.stringify({
            })
        })
        .then(res => res.json())
        .then(data => {

        })
    }

    let $searchContainer = $("#searchContainer")
    function updateResults() {

    }


    $("#searchButton").on("click", function() {
        
    })
})