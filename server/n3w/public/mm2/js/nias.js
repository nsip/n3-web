// nias.js

// global container 
var $n3data = {};

// initiialise selector controls
$(document).ready(function() {
    M.AutoInit();

    $('.modal').modal();


    // setup tab styling on popup
    $('.tabs').tabs();
    // TAB Color
    $(".tabs").css("background-color", '#ffffff00');
    // TAB Indicator/Underline Color
    $(".tabs>.indicator").css("background-color", '#00695c');
    // TAB Text Color
    $(".tabs>li>a").css("color", '#00695c');

    fillTeacherSelector();

    // 
    // respond to teacher selection
    // get the teaching groups
    // 
    $("#teacher-select").on("change", function(event) {
        $selection = $("#teacher-select").val()
        fillTeachingGroupSelector($selection);
        // remove any existing data
        var $progress = $("#progress-table").empty();
        // asssign to var for access in progress logic
        $selectedTeacher = $selection;
    });
    // 
    // respond to the teaching group selection
    // get the prgress data
    // 
    $("#teaching-group-select").on("change", function(event) {
        $selection = $("#teaching-group-select").val()
        fillProgressTable($selection);
    });

});

// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initProgressTableHandler() {

    $('body').off("click"); // clear any existing handlers
    $('body').on('click', '.showmodal tr', function(event) {
        event.preventDefault();
        // alert('clicked');
        fillSyllabusView();
        fillCoverageView();
        showExtendedData();

    });


}


// show the extended data modal form
function showExtendedData() {
    $("#ed-modal").css("max-height", "80%");
    $('#ed-modal').modal('open');
    // re-init tab styles as have been hidden
    $('.tabs').tabs();
    $(".tabs>.indicator").css("background-color", '#00695c');
}



