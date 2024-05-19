package storage

import (
	"time"
)

type MetricsInterface interface {
	ObserveQueryDuration(timeSince time.Time, name string)
}

// TODO Here I used to take repo interface as well, because often I add CacheDecorator and maybe some more and
// Here I didn't use interface just to not to copy it from `Service` package
type MetricDecorator struct {
	db      *InMemoryPostRepository
	metrics MetricsInterface
}

func NewStorageMetricDecorator(db *InMemoryPostRepository, metrics MetricsInterface) *MetricDecorator {
	return &MetricDecorator{
		db:      db,
		metrics: metrics,
	}
}

func (d *MetricDecorator) Create(post Post) error {
	startTime := time.Now()
	err := d.db.Create(post)

	d.metrics.ObserveQueryDuration(startTime, "Create")

	return err
}

func (d *MetricDecorator) GetAll() ([]Post, error) {
	startTime := time.Now()
	posts, err := d.db.GetAll()

	d.metrics.ObserveQueryDuration(startTime, "GetAll")

	return posts, err
}

func (d *MetricDecorator) GetByID(id int) (Post, error) {
	startTime := time.Now()
	post, err := d.db.GetByID(id)

	d.metrics.ObserveQueryDuration(startTime, "GetByID")

	return post, err
}

func (d *MetricDecorator) Update(post Post) error {
	startTime := time.Now()
	err := d.db.Update(post)

	d.metrics.ObserveQueryDuration(startTime, "Update")

	return err
}

func (d *MetricDecorator) Delete(id int) error {
	startTime := time.Now()
	err := d.db.Delete(id)

	d.metrics.ObserveQueryDuration(startTime, "Delete")

	return err
}
