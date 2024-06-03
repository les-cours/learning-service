package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/utils"
)

func (s *Server) CreateClassRooms(ctx context.Context, in *learning.CreateClassRoomsRequest) (*learning.OperationStatus, error) {
	for _, subjectID := range in.SubjectIDs {
		_, err := s.CreateClassRoom(ctx, &learning.CreateClassRoomRequest{
			TeacherID: in.TeacherID,
			SubjectID: subjectID,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
	}
	return &learning.OperationStatus{Success: true}, nil
}

func (s *Server) CreateClassRoom(ctx context.Context, in *learning.CreateClassRoomRequest) (*learning.ClassRoom, error) {

	classRoomID := utils.GenerateUUIDString()
	/*
		GENERATE NAME :
	*/

	var subjectName string
	var arabicSubjectName string

	err := s.DB.QueryRow(`SELECT title ,title_ar FROM subjects 
                         WHERE subject_id = $1 
`, in.SubjectID).Scan(&subjectName, &arabicSubjectName)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("subject")
	}

	var classRoomName = subjectName
	var classRoomArabicName = arabicSubjectName

	/*
		INSERT INTO CLASSROOM THE NEW ONE ...
	*/
	_, err = s.DB.Exec(`INSERT INTO 
    classrooms 
    (classroom_id, teacher_id, subject_id, arabic_title,title,price) 
     VALUES ($1,$2,$3,$4,$5,$6)`, classRoomID, in.TeacherID, in.SubjectID, classRoomArabicName, classRoomName, 1500)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.ClassRoom{
		ClassRoomID:  classRoomID,
		Title:        classRoomName,
		Price:        1500,
		StudentCount: 0,
		Rating:       0,
	}, nil
}

func (s *Server) GetClassRooms(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id

	//teacher_id uuid NULL,
	//subject_id character varying(40) NULL,

	rows, err := s.DB.Query(`
SELECT classroom_id, title, image, price, c.badge, c.description, arabic_title,c.description_ar,t.teacher_id,t.firstname,t.lastname
	FROM classrooms as c
	INNER JOIN public.teachers t on t.teacher_id = c.teacher_id
WHERE c.deleted_at IS NULL AND c.subject_id = $1 ;
        `, subjectID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var classRooms = new(learning.ClassRooms)
	var studentsCount int32
	for rows.Next() {

		var classRoomID, title, image, badge, description, arabicTitle, arabicDescription, teacherID, firstname, lastname string
		var price int32

		err = rows.Scan(&classRoomID, &title, &image, &price, &badge, &description, &arabicTitle, &arabicDescription, &teacherID, &firstname, &lastname)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		classRoom := &learning.ClassRoom{
			ClassRoomID:       classRoomID,
			Title:             title,
			Image:             image,
			Price:             price,
			Badge:             badge,
			Description:       description,
			ArabicTitle:       arabicTitle,
			ArabicDescription: arabicDescription,
			Teacher: &learning.Teacher{
				TeacherID: teacherID,
				Firstname: firstname,
				Lastname:  lastname,
			},
		}

		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

// GetClassRoom acting  like promo
func (s *Server) GetClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.ClassRoom, error) {

	var classRoomID = in.Id
	var studentsCount int32
	var title, image, badge, description, arabicTitle, arabicDescription, teacherID, firstname, lastname string
	var price int32
	var err error

	err = s.DB.QueryRow(`
SELECT  title, image, price, c.badge, c.description, arabic_title,c.description_ar,t.teacher_id,t.firstname,t.lastname
	FROM classrooms as c
	INNER JOIN public.teachers t on t.teacher_id = c.teacher_id
WHERE c.deleted_at IS NULL AND c.classroom_id = $1 ;
        `, classRoomID).Scan(&title, &image, &price, &badge, &description, &arabicTitle, &arabicDescription, &teacherID, &firstname, &lastname)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	classRoom := &learning.ClassRoom{
		ClassRoomID:       classRoomID,
		Title:             title,
		Image:             image,
		Price:             price,
		Badge:             badge,
		Description:       description,
		ArabicTitle:       arabicTitle,
		ArabicDescription: arabicDescription,
		Teacher: &learning.Teacher{
			TeacherID: teacherID,
			Firstname: firstname,
			Lastname:  lastname,
		},
	}

	/*
		Get chapters
	*/

	chapters, err := s.GetChaptersByClassRoom(ctx, &learning.IDRequest{
		Id: classRoomID,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	classRoom.Chapters = chapters
	//get students count for this
	err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	classRoom.StudentCount = studentsCount
	//get Rating  //To Do

	classRoom.Rating = 4.7

	return classRoom, nil
}

func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {
	err := userHasClassRoom(s.DB, in.UserID, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var classRoomID = in.Id

	_, err = s.DB.Exec(`
UPDATE 
    classrooms 
SET deleted_at = $1 
WHERE classroom_id = $1;
        `, classRoomID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.OperationStatus{
		Success: true,
	}, nil
}

func (s *Server) DeleteClassRoomsByTeacher(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	_, err := s.DB.Exec(`
UPDATE 
    classrooms 
SET deleted_at = $1 
WHERE teacher_id = $1;
        `, in.Id)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return &learning.OperationStatus{
		Success: true,
	}, nil
}

func (s *Server) AddDocumentToClassroom(ctx context.Context, in *learning.AddDocumentToClassroomRequest) (*learning.OperationStatus, error) {

	return nil, nil
}

func (s *Server) GetClassRoomsByTeacher(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var teacherID = in.Id

	rows, err := s.DB.Query(`
SELECT classroom_id, title, image, price, badge, description, arabic_title
	FROM classrooms
WHERE deleted_at IS NULL AND teacher_id = $1;
        `, teacherID)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("classroom")
		}
		return nil, ErrInternal
	}
	classRooms := new(learning.ClassRooms)
	var studentsCount int32
	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge, &classRoom.Description, &classRoom.ArabicTitle)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7

		chapters, err := s.GetChaptersByClassRoom(ctx, &learning.IDRequest{
			Id: classRoom.ClassRoomID,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		classRoom.Chapters = chapters

		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

func (s *Server) GetClassRoomsByStudent(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id
	var studentID = in.UserID

	rows, err := s.DB.Query(`
SELECT classrooms.classroom_id, title, image, price, badge, description, arabic_title
	FROM classrooms
	INNER JOIN subscription s on classrooms.classroom_id = s.classroom_id
WHERE deleted_at IS NULL AND subject_id = $1 AND s.student_id = $2;
        `, subjectID, studentID)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("classroom")
		}
		return nil, ErrInternal
	}

	classRooms := &learning.ClassRooms{}
	var studentsCount int32
	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

//func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.GetClassRoomsByStudentRequest) (*learning.ClassRooms, error) {
//
// 1 delete classrooms
//  2 delete videos and pdf from storage ...
//
//	return nil, nil
//}

/*
STUDENTS
*/

func (s *Server) GetMyClassRooms(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id
	var studentID = in.UserID

	rows, err := s.DB.Query(`
SELECT DISTINCT classrooms.classroom_id, title, COALESCE(image, 'default.jpg') AS image,price,COALESCE(badge, '') AS badge, description, arabic_title
	FROM classrooms
	INNER JOIN subscription s on classrooms.classroom_id = s.classroom_id
WHERE s.student_id = $1 AND deleted_at IS NULL AND classrooms.subject_id = $2 ;
        `, studentID, subjectID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	classRooms := &learning.ClassRooms{}
	var studentsCount int32
	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge, &classRoom.Description, &classRoom.ArabicTitle)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount

		classRoom.Rating = 4.7
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

func userHasClassRoom(db *sql.DB, userID, classRoomID string) error {
	var has bool
	err := db.QueryRow(`SELECT EXISTS (
	    SELECT 1
	    FROM classrooms
	    WHERE teacher_id = $1
	    AND classroom_id = $2
	);`, userID, classRoomID).Scan(&has)

	if err != nil {
		return ErrInternal
	}

	if has {
		return nil
	}

	return ErrPermission
}
