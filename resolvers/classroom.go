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

	err := s.DB.QueryRow(`SELECT name ,arabic_name FROM subjects 
                         WHERE subject_id = $1 
`, in.SubjectID).Scan(&subjectName, &arabicSubjectName)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("subject")
	}

	rows, err := s.DB.Query(`SELECT g.name ,g.arabic_name ,g.department_id FROM grades_subjects 
                         inner join public.grades g on g.grade_id = grades_subjects.grade_id
                         WHERE subject_id = $1 
`, in.SubjectID)

	var arabicGradName string
	var gradName string
	var dep string
	for rows.Next() {
		err = rows.Scan(&gradName, &arabicGradName, &dep)
		if err != nil {
			s.Logger.Error(err.Error())

			return nil, ErrInternal
		}
		gradName += " / " + gradName + "( " + dep + " )"
		arabicGradName += " / " + arabicGradName + "( " + dep + " )"
	}

	var classRoomName = subjectName + gradName
	var classRoomArabicName = arabicSubjectName + arabicSubjectName

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

func (s *Server) GetClassRoomsBySubject(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id

	rows, err := s.DB.Query(`
SELECT classrooms.classroom_id,classrooms.title,image,price,badge
FROM classrooms 
    inner join classrooms_document d 
        on classrooms.classroom_id = d.classroom_id
WHERE deleted_at IS NULL AND subject_id = $1;
        `, subjectID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	classRoom := &learning.ClassRoom{}
	classRooms := &learning.ClassRooms{}
	var studentsCount int32
	for rows.Next() {
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

func (s *Server) GetClassRoom(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.ClassRoom, error) {

	var classRoomID = in.Id
	classRoom := &learning.ClassRoom{}
	var studentsCount int32

	err := s.DB.QueryRow(`
SELECT classrooms.classroom_id,classrooms.title,image,price,badge
FROM classrooms 
WHERE deleted_at IS NULL AND classroom_id = $1;
        `, classRoomID).Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge)

	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("classRooms")
		}
		return nil, err
	}
	//get students count for this
	err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoomID).Scan(&studentsCount)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	classRoom.StudentCount = studentsCount
	//get Rating  //To Do
	classRoom.Rating = 4.7

	//get lessons :

	res, err := s.GetLessonsByClassRoom(ctx, &learning.GetByUUIDRequest{
		Id: classRoomID,
	})
	classRoom.Lessons = &learning.Lessons{
		Lessons: res.Lessons,
	}

	return classRoom, nil
}

func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.OperationStatus, error) {

	var classRoomID = in.Id

	_, err := s.DB.Exec(`
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

func (s *Server) AddDocumentToClassroom(ctx context.Context, in *learning.AddDocumentToClassroomRequest) (*learning.OperationStatus, error) {

	return nil, nil
}

//
//func (s *Server) GetClassRoomsByStudent(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.ClassRooms, error) {
//
//	return nil, nil
//}
//
//func (s *Server) GetClassRoomsByTeacher(ctx context.Context, in *learning.GetByUUIDRequest) (*learning.ClassRooms, error) {
//
//	return nil, nil
//}

//func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.GetClassRoomsByStudentRequest) (*learning.ClassRooms, error) {
//
// 1 delete classrooms
//  2 delete videos and pdf from storage ...
//
//	return nil, nil
//}
