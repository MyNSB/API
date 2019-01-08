package Remove

import (
	"Util"
	"errors"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"os"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"Sessions/Parse"
)




// This function takes a sport and removes its directory and data in the json file
/* PARAMS:
	@Sport string compulsory (soccer...)
	@Category bool compulsory (rec/grade)
	@ageCategory string optional
	@Team string optional (A...) - disallowed if ageCategory is not provided
*/
func remove(sport string, category bool, ageCategory string, team string) error {
	// Determine what type of delete is required
	/*
		type 1:
			Delete entire sport
		type 2:
			( Only if grade sport )
			Delete age category
		type 3:
			( Only if grade sport )
			Delete team given age category
	 */


	 // Determine the type of request
	 var reqType int

	 switch {
	 case Util.CompoundIsset(sport, ageCategory, team) && category /* Check that it is a grade sport */:
	 	reqType = 3
	 case Util.CompoundIsset(sport, ageCategory) && category /* Check that it is a grade sport */:
	 	reqType = 2
	 case Util.Isset(sport):
	 	reqType = 1
	 default:
		 return errors.New("sport is not set, invalid input")

	}

	// Determine category
	var categoryString string
	if category {
		categoryString = "Grade"
	} else {
		categoryString = "Rec"
	}

	// Start processing the request

	// No matter what request is chosen the all sports file must be read
	data, _ := ioutil.ReadFile("src/SportLocations/Data/allSports.json")
	// Decode this into the interface we need

	// Start performing the correct request on this file
	switch {
	case reqType == 1:
		// Remove from json
		data = jsonparser.Delete(data, "Sports", categoryString, sport)
		// Remove from file structure
		os.Remove("src/SportLocations/Data/Locations" + categoryString + "/" + sport)
	case reqType == 2:
		data = jsonparser.Delete(data, "Sports", categoryString, sport, "Categories", ageCategory)
		os.Remove("src/SportLocations/Data/Locations" + categoryString + "/" + sport + "/" + ageCategory)
	case reqType == 3:
		data = jsonparser.Delete(data, "Sports", categoryString, sport, "Categories", ageCategory, team)
		os.Remove("src/SportLocations/Data/Locations" + categoryString + "/" + sport + "/" + ageCategory + "/" + team)
	}


	// Set the current allSports.json file data to the one just constructed
	ioutil.WriteFile("src/SportLocations/Data/allSports.json", data, 0644)

	// Return nil indicating that it was successful
	return nil
}





// Handler for delete requests
func DeleteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {



	// Get user details from sessions and determine if they are an admin
	user, err := Parse.ParseSessions(r, w)
	if err != nil || !Util.IsAdmin(user) || !Util.ExistsString(user.Permissions, "sportsDep") {
		Util.Error(403, "Forbidden", "User not logged in or does not have sufficient privileges", "Error!", w)
		return
	}

	r.ParseForm()
	// First convert the req type to a boolean
	reqTypeString := r.Form.Get("Sport_Type")
	var reqType bool
	switch {
	case reqTypeString == "Grade":
		reqType = true
	case reqTypeString == "Rec":
		reqType = false
	default:
		Util.Error(400, "Malformed Request", "Invalid sport type", "Error!", w)
		return
	}



	err = remove(r.Form.Get("Sport_Name"), reqType, r.Form.Get("Age_Category"), r.Form.Get("Team"))
	if err != nil {
		Util.Error(400, "An error occured", err.Error(), "Error!", w)
	}
}