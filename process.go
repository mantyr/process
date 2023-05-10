package process

import (
	"context"
	"errors"
	"os"
	"sync"
)

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type Process interface {
	// SetContext устанавливает контекст процесса, используется для метода Run
	SetContext(config Context) Process

	// Run отправляет заявку на запуск процесса и не дожидается запуска
	// Что бы дождаться запуска процесса воспользуйтесь дополнительными средствами синхранизации
	Run(ctx context.Context) error

	// RunWithContext отправляет заявку на запуск процесса и не дожидается запуска
	// Игнорирует ранее установленный контекст и использует контекст который передан при вызове данного метода
	// Что бы дождаться запуска процесса воспользуйтесь дополнительными средствами синхранизации
	RunWithContext(ctx context.Context, config Context) error

	// Stop отправляет заявку на завершение процесса и не дожидается завершения
	// Что бы дождаться завершения используйте связку Stop и Done
	Stop(ctx context.Context) error

	// Status возвращает статус процесса
	Status() Status

	// Done возвращает канал который закрывается после завершения процесса
	// Если процесс не запущен то возвращается закрытый канал
	Done() <-chan struct{}
}

func NewProcess() Process {
	return &process{
		status: NotRunning,
	}
}

type Collector interface {
	SendStatus(externalID string, status Status)
}

func (p *process) setStatus(status Status) {
	p.status = status
	if p.collector == nil {
		return
	}
	p.collector.SendStatus(p.externalID, status)
}

type process struct {
	mutex sync.RWMutex

	// externalID это внешний идентификатор для процесса
	externalID string

	// collector это внешний механизм для регистрации статуса процесса
	collector Collector

	// context это настройки для запуска процесса операционной системы
	context Context

	// status это статус процесса по принципу state-machine
	status Status

	job struct {
		// cancelFunc это функция для закрытия контекста процесса операционной системы
		cancelFunc context.CancelFunc
	}

	subscribers struct {
		// context это контекст для ожидающих завершения процесса из вне
		context context.Context

		// cancelFunc это функция для закрытия контекста для ожидающих завершения процесса из вне
		cancelFunc context.CancelFunc
	}

	// processPid это идентификатор запущенного процесса или последнего который запускался
	processPid int

	// processState это информация о завершении процессе если таковой запускался
	processState *os.ProcessState
}

const (
	// Up процесс а в процессе запуска
	Up Status = "UP"

	// Running процес запущен и работает
	Running Status = "RUNNING"

	// Down процесс в процессе завершения
	Down Status = "DOWN"

	// NotRunning процесс не запущен
	NotRunning Status = "NOT-RUNNING"
)

type Status string

func (p *process) SetExternalID(externalID string) Process {
	p.externalID = externalID
	return p
}

func (p *process) SetCollector(collector Collector) Process {
	p.collector = collector
	return p
}

func (p *process) SetContext(config Context) Process {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.context = config
	return p
}

func (p *process) Run(ctx context.Context) error {
	return p.RunWithContext(ctx, p.context)
}

// Run запускает процесс ассинхронно, что бы понять удалось ли запустить процесс или нет - воспользуйтесь Wait и Status
func (p *process) RunWithContext(ctx context.Context, config Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch p.status {
	case Up:
		return nil
	case Running:
		return nil
	case Down:
		return errors.New("down")
		// можно дождаться завершения и запустить повторно
	case NotRunning:
		jobContext, cancelFunc := context.WithCancel(ctx)
		p.setStatus(Up)
		p.job.cancelFunc = cancelFunc
		p.subscribers.context, p.subscribers.cancelFunc = context.WithCancel(context.Background())
		go p.run(jobContext, config)
		return nil
	default:
		return errors.New("unexpected process status")
	}
}

func (p *process) Status() Status {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.status
}

func (p *process) Done() <-chan struct{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	switch p.status {
	case Up:
		return p.subscribers.context.Done()
	case Running:
		return p.subscribers.context.Done()
	case Down:
		return p.subscribers.context.Done()
	case NotRunning:
		return closedchan
	default:
		return closedchan
	}
}

func (p *process) Stop(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch p.status {
	case Up:
		p.setStatus(Down)
		go p.job.cancelFunc()
		return nil
	case Running:
		p.setStatus(Down)
		go p.job.cancelFunc()
		return nil
	case Down:
		return nil
	case NotRunning:
		return nil
	default:
		// вернуть ошибку "неизвестный статус процесса"
		return nil
	}
}

func (p *process) run(ctx context.Context, config Context) {
	defer p.subscribers.cancelFunc()
	select {
	case <-ctx.Done():
		p.mutex.Lock()
		p.setStatus(NotRunning)
		p.mutex.Unlock()
		return
	default:
	}

	process, err := config.StartProcess()

	p.mutex.Lock()
	{
		if err != nil {
			p.setStatus(NotRunning)
		} else {
			p.setStatus(Running)
			p.processPid = process.Pid
		}
	}
	p.mutex.Unlock()

	if err != nil {
		return
	}
	go wait(process, ctx)

	state, err := process.Wait()
	if err != nil {
		p.mutex.Lock()
		p.setStatus(NotRunning)
		p.mutex.Unlock()
		return
	}
	p.mutex.Lock()
	p.setStatus(NotRunning)
	p.processState = state
	p.mutex.Unlock()
	return
}

func wait(process *os.Process, ctx context.Context) {
	<-ctx.Done()
	// тут можно дать какой-то сигнал приложению и подождать завершения, если завершение не произошло то вызвать Kill
	process.Kill()
}
