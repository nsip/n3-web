// n3w.js



$(document).ready(function() {
	$("#query").empty()
	$("#filters").empty()
    fillQuery()
    fillFilters()
    $("#btn-search").click(function() {
        // alert("search button clicked");
        $("#staff-personal").empty()
        $("#teaching-group").empty()
        $("#grading-assignment").empty()
        $("#xapi-statement").empty()
        $("#subject").empty()
        $("#syllabus").empty()
        $("#lesson").empty()
        $.ajax({
            method: "GET",
            url: "http://localhost:1323/staffTraversal?staffid\=A4F0069E-D3B8-4822-BDD9-4D649E2A47FD",
            contentType: "application/json"
        }).done(function(result) {
            console.log(result)
            var str = JSON.stringify(result.StaffPersonal, null, 2)
            $("#staff-personal").append("<div class='card-panel teal lighten-4'><pre>" +
                str +
                "</pre></div>");
            var str = JSON.stringify(result.TeachingGroup, null, 2)
            $("#teaching-group").append("<div class='card-panel blue lighten-4'><pre>" +
                str +
                "</pre></div>");
            fillGradingAssignments(result.GradingAssignment)
            fillXAPI(result.XAPI)
            fillSubject(result.Subject)
            fillSyllabus(result.Syllabus)
            fillLessons(result.Lesson)
        });
    });
});

function fillQuery() {
    var query = {
        TraversalSpec: [
            "StaffPersonal",
            "TeachingGroup",
            "GradingAssignment",
            "XAPI",
            "Subject",
            "Syllabus",
            "Lesson"
        ]
    }
    var str = JSON.stringify(query, null, 2)
    $("#query").append("<div class='card-panel grey lighten-4'>" +
        "<h7>Traversal</h7><pre>" +
        str +
        "</pre></div>");
}

function fillFilters() {
    filter = {
		XAPI:[{
			Predicate:"actor.name",TargetValue:"Albert Lombardi"
		}],
		TeachingGroup:[{
			Predicate:".LocalId",TargetValue:"2018-History-8-1-A"
		}]
	}
    var str = JSON.stringify(filter, null, 2)
    $("#filters").append("<div class='card-panel grey lighten-4'>" +
    	"<h7>Filters</h7><pre>" +
        str +
        "</pre></div>");
}

function fillGradingAssignments(data) {
    $.each(data, function() {
        // console.log(this);
        var str = JSON.stringify(this, null, 2)
        $("#grading-assignment").append("<div class='card-panel orange lighten-4'><pre>" +
            str +
            "</pre></div>");
    });
}

function fillXAPI(data) {
    $.each(data, function() {
        // console.log(this);
        var str = JSON.stringify(this, null, 2)
        $("#xapi-statement").append("<div class='card-panel green lighten-4'><pre>" +
            str +
            "</pre></div>");
    });
}

function fillSubject(data) {
    $.each(data, function() {
        // console.log(this);
        var str = JSON.stringify(this, null, 2)
        $("#subject").append("<div class='card-panel teal lighten-4'><pre>" +
            str +
            "</pre></div>");
    });
}

function fillSyllabus(data) {
    $.each(data, function() {
        // console.log(this);
        var str = JSON.stringify(this, null, 2)
        $("#syllabus").append("<div class='truncate card-panel blue lighten-4'><pre>" +
            str +
            "</pre></div>");
    });
}

function fillLessons(data) {
    $.each(data, function() {
        // console.log(this);
        var str = JSON.stringify(this, null, 2)
        $("#lesson").append("<div class='card-panel orange lighten-4'><pre>" +
            str +
            "</pre></div>");
    });
}