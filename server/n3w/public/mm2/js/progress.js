// progress2.js


var $selectedTeacher = "";
var $n3data = {};

// 
// fill the progress data table
// 
function fillProgressTable(tgroup_id) {

    var $progress = $("#progress-table");
    $($progress).empty();

    $.ajax({
        method: "POST",
        url: "http://localhost:1323/n3/graphql",
        headers: { 'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" },
        contentType: "application/json",
        data: JSON.stringify({
            query: progressQuery,
            variables:

            {
                qspec: {
                    queryType: "traversalWithId",
                    queryValue: $selectedTeacher,
                    traversal: [
                        "StaffPersonal",
                        "TeachingGroup",
                        "StudentPersonal",
                        "StudentAttendanceTimeList",
                        "StudentPersonal",
                        "TeachingGroup",
                        "GradingAssignment",
                        "Property.Link",
                        "XAPI",
                        "Property.Link",
                        "Subject",
                        "Unique.Link",
                        "Syllabus",
                        "Unique.Link",
                        "Lesson"
                    ],
                    filters: [
                        { eq: ["TeachingGroup", ".RefId", tgroup_id] }
                    ]
                }
            }
        })
    }).done(function(result) {

        $n3data = result.data.q

        // console.log($n3data);


        // restructure data for display...
        $n3data.assignmentResults = extractAssignmentResults($n3data);
        $n3data.userIds = extractUserIds($n3data);
        $n3data.attendanceReports = extractAttendanceReports($n3data);

        $($progress).append("<h5>" + $n3data.TeachingGroup[0].ShortName + "</h5>");
        for (const gA of $n3data.GradingAssignment.sort(assignmentCompare)) {
            $($progress).append("<h5>" + gA.Description + "</h5>");
            // add a results table
            var $table = $('<table>').addClass('highlight showmodal');
            // add table header
            $table.append("<thead><tr><th>Name</th><th>Success</th><th>Completed</th><th>Score</th><th>Time On Task</th><th>Days Absent</th></tr></thead>");
            var $studentResults = $('<tbody>');
            for (const xapi of $n3data.assignmentResults[gA.DetailedDescriptionURL].sort(xapiCompare)) {
                var $resultsRow = $('<tr>');
                $resultsRow.append("<td>" + xapi.actor.name + "</td>");
                $resultsRow.append("<td>" + xapi.result.success + "</td>");
                $resultsRow.append("<td>" + xapi.result.completion + "</td>");
                $resultsRow.append("<td>" + xapi.result.score.scaled + "</td>");
                $resultsRow.append("<td>" + minutesFromIsoDuration(xapi.result.duration) + " minutes</td>");
                // fetch the attendance data
                refid = $n3data.userIds[xapi.actor.mbox];
                attArray = $n3data.attendanceReports[refid];
                if (!attArray || !attArray.length) {
                    attArray = [];
                }
                $resultsRow.append("<td>" + attArray.length + "</td>");
                $studentResults.append($resultsRow);
                $table.append($studentResults);
            }
            $progress.append($table);
        }
    });
    initProgressTableHandler();

}


// 
// sort grading assignemts by description
// 
function assignmentCompare(a, b) {
    // Use toUpperCase() to ignore character casing
    const gaA = a.Description.toUpperCase();
    const gaB = b.Description.toUpperCase();

    let comparison = 0;
    if (gaA > gaB) {
        comparison = 1;
    } else if (gaA < gaB) {
        comparison = -1;
    }
    return comparison;
}

// 
// sort xapi results by score
// 
function xapiCompare(a, b) {
    const gaA = a.result.score.scaled;
    const gaB = b.result.score.scaled;

    let comparison = 0;
    if (gaA > gaB) {
        comparison = 1;
    } else if (gaA < gaB) {
        comparison = -1;
    }
    return comparison * -1;
}



// 
// groups xapi results by grading assignment
// 
function extractAssignmentResults(queryData) {

    results = {};

    for (const xapi of queryData.XAPI) {
        // console.log(xapi);
        xapiArray = results[xapi.object.id];
        if (!xapiArray || !xapiArray.length) {
            results[xapi.object.id] = [];
        }
        results[xapi.object.id].push(xapi);
    }

    return results;
}

// 
// access sif refids via xapi mbox
// 
function extractUserIds(queryData) {

    results = {};

    for (const sp of queryData.StudentPersonal) {
        // console.log(sp);
        mbox = sp.PersonInfo.EmailList.Email[0].value;
        spArray = results[mbox];
        if (!spArray || !spArray.length) {
            results[mbox] = [];
        }
        results[mbox].push(sp.RefId);
    }

    return results;
}

// 
// group attendance by student
// 
function extractAttendanceReports(queryData) {

    results = {};

    for (const satl of queryData.StudentAttendanceTimeList) {
        // console.log(sp);
        sprefid = satl.StudentPersonalRefId;
        satlArray = results[sprefid];
        if (!satlArray || !satlArray.length) {
            results[sprefid] = [];
        }
        results[sprefid].push(satl);
    }

    return results;
}



// 
// helpers for parsing ISO durations from xAPI
// 
var regex = /P((([0-9]*\.?[0-9]*)Y)?(([0-9]*\.?[0-9]*)M)?(([0-9]*\.?[0-9]*)W)?(([0-9]*\.?[0-9]*)D)?)?(T(([0-9]*\.?[0-9]*)H)?(([0-9]*\.?[0-9]*)M)?(([0-9]*\.?[0-9]*)S)?)?/

minutesFromIsoDuration = function(duration) {
    var matches = duration.match(regex);

    return parseFloat(matches[14]) || 0;
}


// 
// the GQL query
// 
var progressQuery = `
query fullTraversal($qspec: QueryInput!) {
  q(qspec: $qspec) {
    StudentAttendanceTimeList{
      RefId
      StudentPersonalRefId
      Date
      AttendanceTimes{
        AttendanceTime{
          AttendanceCode {
            Code
          }
          AttendanceStatus
          DurationValue
          EndTime
          StartTime
        }
      }
    }
    Syllabus {
      learning_area
      stage
      subject
      concepts {
        name
        description
      }
      inquiry_skills {
        name
        skills{
          ac
          skill
        }
      }
      tools{
        name
        examples
      }
    }
    StaffPersonal {
      LocalId
      RefId
      EmploymentStatus
    }
    GradingAssignment {
      RefId
      DetailedDescriptionURL
      Description
      PointsPossible
    }
    Subject {
      subject
      learning_area
      stage
      yrLvls
    }
    Lesson {
      lesson_id
      title
    }
    SchoolInfo {
      RefId
    }
    StudentPersonal {
      RefId
      LocalId
      PersonInfo {
        Name {
          FamilyName
          GivenName
        }
        EmailList{
          Email{
            value
          }
        }
      }
    }
    TeachingGroup {
      SchoolYear
      LongName
      ShortName
      RefId
      LocalId
    }
    XAPI {
      id
      actor {
        mbox
        name
      }
      verb {
        id
        display {
          en_US
        }
      }
      object {
        id
        definition {
          name
          type
        }
      }
      result {
        duration
        success
        completion
        score {
          min
          max
          scaled
        }
      }
    }
  }
}
`