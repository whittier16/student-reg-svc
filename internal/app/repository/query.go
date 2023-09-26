package repository

const (
	createStudent = `
	INSERT INTO student (
		email,
		name
	) VALUES (
		?,
		?
	)
`
	getStudentQuery = `
	SELECT
		email,
		name,
		created_on,
		deleted_on,
		updated_on
	FROM
		student
	WHERE
		email=?
`
	getTeacherQuery = `
	SELECT
		email,
		name,
		created_on,
		deleted_on,
		updated_on
	FROM
		teacher
	WHERE
		email=?
`
	createTeacher = `
	INSERT INTO teacher (
		email,
		name
	) VALUES (
		?,
		?
	)
`
	createRegister = `
	INSERT INTO register (
		student_id,
		teacher_id
	) VALUES (
		?,
		?
	)
`
	getRegisteredStudents = `
		SELECT
			student_id
		FROM
			register
		WHERE
			title like ?
		ORDER BY
			title
`
	updatePostQuery = `
	UPDATE
		register
	SET
		suspended_on = NOW()
	WHERE
		student_id = ?
`

	deletePostByID = `
	DELETE FROM
		post
	WHERE
		id = ?
`
)
