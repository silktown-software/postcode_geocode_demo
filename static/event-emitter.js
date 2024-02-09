export class EventEmitter {
    /**
     * Ctor
     */
    constructor() {
        this.funcMap = {};
    }

    /**
     * Registers a handler to an event.
     * @param eventName - the event name as a string
     * @param handler - the handler to execute
     */
    on(eventName, handler) {
        if (!eventName) {
            throw new Error('eventName is undefined')
        }

        if (typeof handler !== 'function') {
            throw new Error('handler is not a function')
        }

        let handlers = this.funcMap[eventName];

        handlers = handlers === undefined ? [handler] : [...handlers, handler];

        this.funcMap[eventName] = handlers
    }

    /**
     * Removes a handler for a specified event
     * @param eventName - the event name as a string
     * @param handler - the handler
     */
    off(eventName, handler) {
        //TODO: implement later on
    }

    /**
     * Fires an event
     * @param eventName - the event to fire
     * @param data - any data to send to the function
     */
    fire(eventName, data = null) {
        if (!eventName) {
            throw new Error('eventName is undefined');
        }

        const funcs = this.funcMap[eventName];

        if (!funcs) {
            return;
        }

        for (const f of funcs) {
            f(data);
        }
    }
}