pages: events.ics;
	mkdir -p public
	cp events.ics public

events.ics: events
	go build -o ical gitlab.com/mirukakoro/sisikyo/events/cmd/ical
	./ical > events.ics
