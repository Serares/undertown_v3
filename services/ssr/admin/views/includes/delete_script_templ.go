// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package includes

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"

import "github.com/Serares/ssr/admin/types"

func HandleDeleteButton(props types.DeleteScriptProps) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_HandleDeleteButton_bb5d`,
		Function: `function __templ_HandleDeleteButton_bb5d(props){document.getElementById('delete_button').addEventListener('click', function(e) {
        e.preventDefault(); // Stop the form from submitting normally
        if(window.confirm("Are you sure you want to delete?")) {
        fetch(props.DeleteUrl, {
            method: 'DELETE', // Set the method to DELETE
            headers: {
                // If you need to send a JSON body with fetch, uncomment and use the lines below
                'Content-Type': 'application/json',
                // Optionally, include headers for CSRF protection as needed
                'X-Requested-With': 'XMLHttpRequest',
                // 'X-CSRF-TOKEN': 'your_csrf_token_here'
            },
            // If sending JSON, stringify your payload and include it in the body
            // body: JSON.stringify({ key: 'value' })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            if (response.redirected) {
                window.location.href = response.url;
                return
            }
            return response.text()
        })
        .then(data => {
            // console.log(data); // Handle success response
            document.documentElement.innerHTML = data;
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
             $("#myModal .modal-body").html("Failed deleting the property")
            $("#myModal").modal("show")
        });
        }
    });
}`,
		Call:       templ.SafeScript(`__templ_HandleDeleteButton_bb5d`, props),
		CallInline: templ.SafeScriptInline(`__templ_HandleDeleteButton_bb5d`, props),
	}
}
