package task

type Task interface {
	Run()
	SetID(id int)
}

type TaskFactory interface {
	CreateTask() Task
}
