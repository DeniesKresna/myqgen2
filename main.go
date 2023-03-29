package main

import (
	"fmt"
	"time"

	"github.com/DeniesKresna/gohelper/utlog"
	"github.com/DeniesKresna/myqgen2/qgen"
)

type User struct {
	ID        int64      `json:"id" db:"id" sqlq:"userID"`
	CreatedBy string     `json:"created_by" db:"created_by" sqlq:"userCreatedBy"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" sqlq:"userCreatedAt"`
	UpdatedBy string     `json:"updated_by" db:"updated_by" sqlq:"userUpdatedBy"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" sqlq:"userUpdatedAt"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by" sqlq:"userDeletedBy"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at" sqlq:"userDeletedAt"`
	FirstName string     `json:"first_name" db:"first_name" sqlq:"userFirstName"`
	LastName  string     `json:"last_name" db:"last_name" sqlq:"userLastName"`
	Email     string     `json:"email" db:"email" sqlq:"userEmail"`
	Phone     string     `json:"phone" db:"phone" sqlq:"userPhone"`
	ImageUrl  *string    `json:"image_url" db:"image_url" sqlq:"userImageURL"`
	Password  string     `json:"-" db:"password" sqlq:"userPassword"`
	Active    int        `db:"active" sqlq:"userActive"`
	RoleId    int64      `json:"role_id" db:"role_id" sqlq:"userRoleID"`
}

func (u User) GetTableName() string {
	return "users"
}

type Role struct {
	ID        int64      `json:"id" db:"id" sqlq:"roleID"`
	CreatedBy string     `json:"created_by" db:"created_by" sqlq:"roleCreatedBy"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" sqlq:"roleCreatedAt"`
	UpdatedBy string     `json:"updated_by" db:"updated_by" sqlq:"roleUpdatedBy"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" sqlq:"roleUpdatedAt"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by" sqlq:"roleDeletedBy"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at" sqlq:"roleDeletedAt"`
	Name      string     `json:"name" db:"name" sqlq:"roleName"`
}

func (u Role) GetTableName() string {
	return "roles"
}

type ExpertEducation struct {
	School    string     `json:"school" db:"education_school" jsondb:"datas>$.education.school"`
	Degree    string     `json:"degree" db:"education_degree" jsondb:"datas>$.education.degree"`
	StartDate *time.Time `json:"start_date" db:"education_start_date" jsondb:"datas>$.education.start_date"`
	EndDate   *time.Time `json:"end_date" db:"education_end_date" jsondb:"datas>$.education.end_date"`
}

type ExpertExperience struct {
	Title     string     `json:"title" db:"experience_title" jsondb:"datas>$.experiences.invite_expert"`
	StartDate *time.Time `json:"start_date" db:"experience_start_date" jsondb:"datas>$.experiences.invite_expert"`
	EndDate   *time.Time `json:"end_date" db:"experience_end_date" jsondb:"datas>$.experiences.invite_expert"`
	Location  string     `json:"location" db:"experience_location" jsondb:"datas>$.experiences.location"`
}

type ExpertConsultation struct {
	DayRecomendations      []string `json:"day_recomendations" db:"consultation_day_recomendations" jsondb:"datas>$.service.consultation.day_recomendations"`
	TimeRecomendationStart int      `json:"time_recomendation_start" db:"consultation_time_recomendation_start" jsondb:"datas>$.service.consultation.invite_expert"`
	TimeRecomendationEnd   int      `json:"time_recomendation_end" db:"consultation_time_recomendation_end" jsondb:"datas>$.service.consultation.time_recomendation_end"`
	MaterialTypes          []string `json:"material_types" db:"consultation_material_types" jsondb:"datas>$.service.consultation.material_types"`
	Fee                    int      `json:"fee" db:"consultation_fee" jsondb:"datas>$.service.consultation.fee"`
}

type ExpertTraining struct {
	DayRecomendation       []string `json:"day_recomendations" db:"training_day_recomendations" jsondb:"datas>$.service.training.day_recomendations"`
	TimeRecomendationStart int      `json:"time_recomendation_start" db:"training_time_recomendation_start" jsondb:"datas>$.service.training.time_recomendation_start"`
	TimeRecomendationEnd   int      `json:"time_recomendation_end" db:"training_time_recomendation_end" jsondb:"datas>$.service.training.time_recomendation_end"`
}

type ExpertInvitation struct {
	DayRecomendation       []string `json:"day_recomendations" db:"invite_day_recomendations" jsondb:"datas>$.service.invite_expert.day_recomendations"`
	TimeRecomendationStart int      `json:"time_recomendation_start" db:"invite_time_recomendation_start" jsondb:"datas>$.service.invite_expert.time_recomendation_start"`
	TimeRecomendationEnd   int      `json:"time_recomendation_end" db:"invite_time_recomendation_end" jsondb:"datas>$.service.invite_expert.time_recomendation_end"`
}

type RecruitExpert struct {
	Capabilities          []string `json:"capabilities" db:"recruit_capabilities" jsondb:"datas>$.service.recruit_expert.capabilities"`
	AcceptableTypeOfWorks []string `json:"acceptable_type_of_works" db:"recruit_acceptable_type_of_works" jsondb:"datas>$.service.recruit_expert.acceptable_type_of_works"`
}

type Expert struct {
	ID            int64              `json:"id" sqlq:"expertID" db:"id" jsondb:"datas>$.id>UNSIGNED"`
	CreatedBy     string             `json:"created_by" sqlq:"expertCreatedBy" db:"created_by" jsondb:"datas>$.created_by"`
	CreatedAt     string             `json:"created_at" sqlq:"expertCreatedAt" db:"created_at" jsondb:"datas>$.created_at"`
	UpdatedBy     string             `json:"updated_by" db:"updated_by" sqlq:"expertUpdatedBy" jsondb:"datas>$.updated_by"`
	UpdatedAt     string             `json:"updated_at" db:"updated_at" sqlq:"expertUpdatedAt" jsondb:"datas>$.updated_at"`
	DeletedBy     *string            `json:"deleted_by" db:"deleted_by" sqlq:"expertDeletedBy" jsondb:"datas>$.deleted_by"`
	DeletedAt     *string            `json:"deleted_at" db:"deleted_at" sqlq:"expertDeletedAt" jsondb:"datas>$.deleted_at>CHAR(20)"`
	Name          string             `json:"name" db:"name" sqlq:"expertName" jsondb:"datas>$.name>CHAR(50)"`
	Description   string             `json:"description" db:"description" sqlq:"expertDescription" jsondb:"datas>$.description"`
	Profession    string             `json:"profession" db:"profession" sqlq:"expertProfession" jsondb:"datas>$.profession"`
	Company       string             `json:"company" db:"company" sqlq:"expertCompany" jsondb:"datas>$.company"`
	Domicile      string             `json:"domicile" db:"domicile" db:"domicile" sqlq:"expertDomicile" jsondb:"datas>$.domicile"`
	Education     ExpertEducation    `json:"education" db:"education" sqlq:"expertEducation" jsondb:"datas>$.education"`
	Experiences   []ExpertExperience `json:"experiences" db:"experiences" sqlq:"expertExperiences" jsondb:"datas>$.experiences"`
	ExperienceYOE int                `json:"experience_yoe" db:"experience_yoe" sqlq:"expertExperienceYOE" jsondb:"datas>$.experience_yoe"`
	SocialMedia   struct {
		Linkedin string `json:"linkedin" db:"linkedin" sqlq:"expertLinkedin" jsondb:"datas>$.social_media.linkedin"`
		Facebook string `json:"facebook" db:"facebook" sqlq:"expertFacebook" jsondb:"datas>$.social_media.facebook"`
	} `json:"social_media" db:"social_media" sqlq:"expertSocialMedia" jsondb:"datas>$.social_media"`
	Sectors           []string `json:"sectors" db:"sectors" sqlq:"expertSectors" jsondb:"datas>$.sectors"`
	AvailableServices []string `json:"available_services" db:"available_services" sqlq:"expertAvailableServices" jsondb:"datas>$.available_services"`
	Service           struct {
		Consultation ExpertConsultation `json:"consultation" db:"consultation" sqlq:"expertConsultation" jsondb:"datas>$.service.consultation"`
		Training     ExpertTraining     `json:"training" db:"training" sqlq:"expertTraining" jsondb:"datas>$.service.training"`
		Invitation   ExpertInvitation   `json:"invite_expert" db:"invite_expert" sqlq:"expertInviteExpert" jsondb:"datas>$.service.invite_expert"`
		Recruitment  RecruitExpert      `json:"recruit_expert" db:"recruit_expert" sqlq:"expertRecruitExpert" jsondb:"datas>$.service.recruit_expert"`
	} `json:"service" db:"service" sqlq:"expertService" jsondb:"datas>$.service"`
	Star   float32 `json:"star" db:"star" sqlq:"expertStar" jsondb:"datas>$.star"`
	Image  string  `json:"image" db:"image" sqlq:"expertImage" jsondb:"datas>$.image"`
	Active int     `json:"active" db:"active" sqlq:"expertActive" jsondb:"datas>$.active"`
	Datas  string  `json:"-" db:"datas" sqlq:"expertDatas"`
	UserID string  `json:"-" db:"user_id" sqlq:"expertUserID"`
}

func (e Expert) GetTableName() string {
	return "experts"
}

func main() {
	qGenObj, err := qgen.InitObject(true, User{}, Role{}, Expert{})
	if err != nil {
		utlog.Errorf("error: %+v", err)
		return
	}

	// query := `
	// 	{
	// 		"select": [
	// 			{"col": "u.*"},
	// 			{"col": "r.name", "as": "role_name", "value": "r.name"}
	// 		],
	// 		"from": {
	// 			"value": "users", "as": "u"
	// 		},
	// 		"join": [
	// 			{"value": "roles", "as": "r", "type": "inner", "conn": "r.id = u.role_id"}
	// 		],
	// 		"where": {
	// 			"and": [
	// 				{"col":"id", "value":"u.id"},
	// 				{"col":"email", "value":"u.email"},
	// 				{"col":"-", "value":"u.active = 1"}
	// 			]
	// 		}
	//   	}
	// `

	query := `
		{
			"select": [
				{"col": "u.*"},
				{"col": "r.name", "as": "role_name", "value": "r.name"}
			],
			"from": {
				"value": "users", "as": "u"
			},
			"join": [
				{"value": "roles", "as": "r", "type": "inner", "conn": "r.id = u.role_id"}
			],
			"where": {
				"and": [
					{"col":"ids", "value":{
						"select": [
							{"col": "-", "value": "u2.id"}
						],
						"from": {
							"value": "users", "as": "u2"
						}
					}},
					{"col":"email", "value":"u.email"},
					{"col":"-", "value":"u.active = 1"}
				]
			}
	  	}
	`

	// query := `
	// 	{
	// 		"select": [
	// 			{"col": "e.*"}
	// 		],
	// 		"from": {
	// 			"value": "experts", "as": "e"
	// 		},
	// 		"where": {
	// 			"and": [
	// 				{"col":"name", "value":"e.name"}
	// 			]
	// 		}
	//   	}
	// `

	// jsondb:"datas>$.name"

	args := qgen.Args{
		Fields: []string{
			"u.id",
			"u.first_name",
		},
		Conditions: map[string]interface{}{
			"email:like": "%info%",
		},
		Sorting: []string{"-id"},
		Limit:   3,
	}

	res := qGenObj.Build(query, args)
	fmt.Println(res)
}
