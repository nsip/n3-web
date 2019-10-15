// coverage.js


// 
// shows the lesson sequences seen by the selected student
// 
function fillCoverageView() {

    var $sequences = $("#sequences");
    $($sequences).empty();

    $("#coverage-status-message").empty();

    // only hsie syllabus supported
    var $selected = $("#teaching-group-select option:selected").text();
    var $learning_area = "";
    if ($selected.indexOf("Geography") != -1) {
        $learning_area = "geography";
    }
    if ($selected.indexOf("History") != -1) {
        $learning_area = "history";
    }

    if (($learning_area != "history") && ($learning_area != "geography")) {
        $("#coverage-status-message").append("<p>Syllabus coverage not available for this Learning Area</p>");
        return
    }

    if (!$n3data.Syllabus) {
        $("#coverage-status-message").append("<p>Syllabus coverage not available for this Learning Area</p>");
        return        
    }

    // if no lesson sequences have been created
    if (!$n3data.LessonSequence) {
        $("#coverage-status-message").append("<p>No Lesson Sequences created for this Learning Area</p>");
        return        
    }

    // add column header
    $sequences.append($("<p><h6>Course Attendance</h6></p>"));
    $.each($n3data.LessonSequence, function() {
        var $lessonBlock = $('<div>');
        $lessonBlock.append($("<p><h6>" +
            this.thecourse +
            " (" + this.thedescription +
            ", taught by: " + this.userId + "), " +
            "</h6></p>"));
        var $list = $('<ul>').addClass("collection");
        $.each(this.lessonList, function() {
            var $listItem = $("<li>").addClass("collection-item");
            $listItem.append(makeCheckbox(this.summary));
            // $listItem.text(this.summary);
            $($list).append($listItem);
        });
        $lessonBlock.append($list);
        $sequences.append($lessonBlock)
    });

}