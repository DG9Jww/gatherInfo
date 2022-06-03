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

	srcPort uint16

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
	size int
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
	if link.isEmpty() {
		link.head = newTab
		link.tail = newTab
		return
	}

	link.tail.next = newTab
	link.tail = newTab
	newTab.pre = link.tail
	newTab.next = link.head
	link.head.pre = newTab
	link.size++

}

//remove statusTable
func (link *tableLinkList) remove(tab *statusTable) {
	link.lock.Lock()
	defer link.lock.Unlock()
	if tab == link.head {
		link.head = link.head.next
	}
	if tab == link.tail {
		link.tail = link.tail.pre
	}
	tab.pre.next = tab.next
	tab.next.pre = tab.pre
	tab = nil
	link.size--
}

//queryStatusTable according to subdomain name and port
func (link *tableLinkList) queryStatusTab(subdomain string, port uint16) (*statusTable, error) {
	if link.head == nil {
		return nil, emptyLink
	}
	current := link.head
	for {
		if current.domain == subdomain && current.srcPort == port {
			return current, nil
		}
		if current == link.tail {
			return nil, notFound
		}
		current = current.next
	}
}

func (link *tableLinkList) printAllTables() {
	current := link.head
	for {
		fmt.Println("node:", *current)
		if current == link.tail {
			break
		}
		current = current.next
	}
}
