// teachinggroups.js


// 
// fill the teaching group selector
// 
function fillTeachingGroupSelector(teacher_guid) {

    var $dropdown = $("#teaching-group-select");
    $($dropdown).empty();


    $.ajax({
        method: "POST",
        url: "http://localhost:1323/n3/graphql",
        contentType: "application/json",
        headers: { 'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" },
        data: JSON.stringify({
            query: tgQuery,
            variables: {
                qspec: {
                    queryType: "traversalWithId",
                    queryValue: teacher_guid,
                    "traversal": [
                        "StaffPersonal",
                        "TeachingGroup"
                    ]
                }
            }
        })
    }).done(function(result) {
        // fill teaching group selector
        $.each(result.data.q.TeachingGroup, function() {
            // get the teaching periods
            var periodDescriptor = ""
            $.each(this.TeachingGroupPeriodList.TeachingGroupPeriod, function() {
                periodDescriptor = periodDescriptor + "(" + this.DayId + ":" + this.PeriodId + ") "
            });
            $dropdown.prepend($("<option />").val(this.RefId).text(this.ShortName + " - " + periodDescriptor));
        });
        // have to re-initialise component to render
        $($dropdown).formSelect();
    });

}


var tgQuery = `query fullTraversal($qspec: QueryInput!) {
  q(qspec: $qspec) {
    TeachingGroup {
      SchoolYear
      LocalId
      LongName
      ShortName
      TimeTableSubjectRefId
      RefId
      TeachingGroupPeriodList {
        TeachingGroupPeriod {
          DayId
          PeriodId
          RoomNumber
        }
      }
    }
  }
}

`