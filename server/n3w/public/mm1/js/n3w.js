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
        // query
        var queryText = staffTraversalQuery();
        $.ajax({
            method: "POST",
            url: "http://localhost:1323/n3/graphql",
            contentType: "application/json",
            headers: { 'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" },
            data: JSON.stringify({
                query: queryText,
                variables: {
                    qspec: {
                        queryType: "traversalWithId",
                        queryValue: "A4F0069E-D3B8-4822-BDD9-4D649E2A47FD",
                        traversal: [
                            "StaffPersonal",
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
                            { eq: ["XAPI", "actor.name", "Albert Lombardi"] },
                            { eq: ["TeachingGroup", ".LocalId", "2018-History-8-1-A"] }
                        ]
                    }
                }
            })
        }).done(function(result) {
            console.log(result.data.q)
            var str = JSON.stringify(result.data.q.StaffPersonal, null, 2)
            $("#staff-personal").append("<div class='card-panel teal lighten-4'><pre>" +
                str +
                "</pre></div>");
            var str = JSON.stringify(result.data.q.TeachingGroup, null, 2)
            $("#teaching-group").append("<div class='card-panel blue lighten-4'><pre>" +
                str +
                "</pre></div>");
            fillGradingAssignments(result.data.q.GradingAssignment)
            fillXAPI(result.data.q.XAPI)
            fillSubject(result.data.q.Subject)
            fillSyllabus(result.data.q.Syllabus)
            fillLessons(result.data.q.Lesson)
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
        "filters": [
            { "eq": ["XAPI", "actor.name", "Albert Lombardi"] },
            { "eq": ["TeachingGroup", ".LocalId", "2018-History-8-1-A"] }
        ]
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





function staffTraversalQuery() {
    return `query fullTraversal($qspec: QueryInput!) {
  q(qspec: $qspec) {
    Syllabus {
      learning_area
      stage
      subject
    courses {
      name
      outcomes {
        id
        description
      }
      lifeskills_outcomes
      inquiry_questions
      focus
      content_areas {
        name
        investigations {
          description
          ac
          examples
        }
      }
    }
        overview
    concepts {
      name
      description
    }
    inquiry_skills {
      name
      skills {
        skill
        ac
      }
    }
    tools {
      name
      examples
      code
    }

    }
    StaffPersonal {
      LocalId
      RefId
      EmploymentStatus
      PersonInfo {
        Name {
            FamilyName
            GivenName
            Type
        }
        Demographics {
            Sex
        }
      }
    }
    GradingAssignment {
      DetailedDescriptionURL
      PointsPossible
      Description
      TeachingGroupRefId
      LearningStandardList {
        LearningStandard {
            LearningStandardLocalId
        }
      }
      RefId
    }
    Subject {
      subject
      learning_area
      stage
      yrLvls
      synonyms
    }
    Lesson {
      lesson_id
      content
      title
      stage
      subject
      teacher
      learning_area
    }
    SchoolInfo {
      StateProvinceId
      SchoolURL
      SchoolType
      RefId
      SchoolDistrict
      LocalId
      SchoolName
      CommonwealthId
      SchoolSector
    }
    StudentPersonal {
      RefId
      LocalId
      PersonInfo {
        Demographics {
          BirthDate
          IndigenousStatus
          Sex
        }
        Name {
            FamilyName
            GivenName
        }
      }
    }
    TeachingGroup {
      SchoolYear
      LocalId
      LongName
      ShortName
      TimeTableSubjectRefId
      TeachingGroupPeriodList {
        TeachingGroupPeriod {
            DayId
            PeriodId
            RoomNumber
        }
      }
      TeacherList {
        TeachingGroupTeacher {
            StudentPersonalRefId
            Association
            Name {
                FamilyName
                GivenName
            }
        }
      }
      StudentList {
        TeachingGroupStudent {
            StudentPersonalRefId
            Name {
                FamilyName
                GivenName
            }
        }
      }
      RefId
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

}