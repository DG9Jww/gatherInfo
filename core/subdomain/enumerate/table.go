package enumerate

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	notFound  = errors.New("query error:can not find specific status table")
	emptyLink = errors.New("status tables is empty")
)

//status table,use doubly circular linklist
type statusTable struct {
	domain string

	flagID uint16

	//sending packet time
	time time.Time

	//last resolver used
	resolver string

	//the amount of attempts
	retry int8

	//status,0 unsent,1 sent
	status uint8

	//doubly circular linklist
	pre  *statusTable
	next *statusTable
}

type tableLinkList struct {
	done int64
	size int64
	head *statusTable
	tail *statusTable
	lock sync.Mutex
}

func initTabLinkList() *tableLinkList {
	return &tableLinkList{}
}

func (link *tableLinkList) isEmpty() bool {
	return link.head == nil
}

//append statusTable
func (link *tableLinkList) append(newTab *statusTable) {
	link.lock.Lock()
	defer link.lock.Unlock()
	link.size++
	if link.isEmpty() {
		link.head = newTab
		link.tail = newTab
		return
	}

	newTab.pre = link.tail
	newTab.next = link.head
	link.tail.next = newTab
	link.head.pre = newTab
	link.tail = newTab

}

//remove statusTable
func (link *tableLinkList) remove(tab *statusTable) error {
	link.lock.Lock()
	defer link.lock.Unlock()
	if link.isEmpty() {
		return emptyLink
	}
	if tab == nil {
		return notFound
	}
	//the last node
	if tab == tab.pre {
		link.tail = nil
		link.head = nil
		tab = nil
		link.done++
		return nil
	}
	if tab == link.head {
		link.head = tab.next
	}
	if tab == link.tail {
		link.tail = tab.pre
	}
	tab.pre.next = tab.next
	tab.next.pre = tab.pre
	tab = nil
	link.done++
	return nil
}

//queryStatusTable according to subdomain name and port
func (link *tableLinkList) queryStatusTab(subdomain string, flagID uint16) (*statusTable, error) {
	if link.head == nil {
		return nil, emptyLink
	}
	current := link.head
	for {
		if current.domain == subdomain && current.flagID == flagID {
			return current, nil
		}
		if current == link.tail {
			return nil, notFound
		}
		current = current.next
	}
}

func (link *tableLinkList) printAllTables() {
	if link.isEmpty() {
		return
	}
	current := link.head
	for {
		fmt.Println("node:", *current)
		if current == link.tail {
			break
		}
		current = current.next
	}
}
