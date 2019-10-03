// teachers.js





// 
// fill the teachers selector
// 
function fillTeacherSelector() {
    var $dropdown = $("#teacher-select");
    $($dropdown).empty();

    $.ajax({
        method: "POST",
        url: "http://localhost:1323/n3/graphql",
        contentType: "application/json",
        headers: { 'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" },
        data: JSON.stringify({
            query: teacherQuery,
            variables: {
                qspec: {
                    queryType: "findByType",
                    queryValue: "StaffPersonal"
                }
            }
        })
    }).done(function(result) {
        // fill the teachers dropdown
        $.each(result.data.q.StaffPersonal, function() {
            // console.log(this);
            $dropdown.prepend($("<option />").val(this.RefId).text(this.PersonInfo.Name.GivenName +
                " " + this.PersonInfo.Name.FamilyName));
        });
        // have to re-initialise component to render
        $($dropdown).formSelect();
    });
}


var teacherQuery = `
query teachersQuery($qspec: QueryInput!) {
  q(qspec: $qspec) {
    StaffPersonal {
      RefId
      PersonInfo {
        Name {
          FamilyName
          GivenName
        }
      }
    }
  }
}
`