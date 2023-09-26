package services

type SendNotificationsParams struct {
	Teacher       string `json:"teacher" valid:"email,required"`
	Notifications string `json:"notifications" valid:"required"`
}

type RegisterStudentsParams struct {
	TeacherEmail  string   `json:"teacher" valid:"email,required"`
	StudentEmails []string `json:"students" valid:"email,required"`
}

type CreateStudentParams struct {
	Email string `json:"email" valid:"email,required"`
	Name  string `json:"name"  valid:"optional"`
}

type CreateTeacherParams struct {
	Email string `json:"email" valid:"email,required"`
	Name  string `json:"name"  valid:"optional"`
}

type SuspendStudentsParams struct {
	Student string `json:"student" valid:"email,required"`
}

type GetCommonStudentsParams struct {
	Teacher []string `json:"teacher" valid:"email,required"`
}
