import {EventEmitter} from './event-emitter.js';

export class PostcodeForm {
    /**
     * Ctor
     * @param name - the form name
     * @param eventEmitter - an instance of the event emitter
     */
    constructor(name, eventEmitter) {
        if (!name) {
            throw new Error('selector string is empty');
        }

        if (eventEmitter && typeof eventEmitter === typeof (EventEmitter)) {
            throw new Error('event emitter is undefined')
        }

        this.form = document.forms[name];
        this.eventEmitter = eventEmitter;
    }

    /**
     * Initialises form DOM elements
     */
    init() {
        this.postcodeTextBox = this.form.postalcode;
        this.submitButton = this.form.querySelector('button[type="submit"]');
        this.resetButton = this.form.querySelector('button[type="reset"]');
        this.errorMessage = this.form.querySelector('.error-message');
    }

    /**
     * Binds the form elements and events
     */
    bind() {
        this.form.addEventListener('submit', (evt) => this.#submitHandler(evt));
        this.postcodeTextBox.addEventListener('keyup', (evt) => this.#handlePostCodeKeyUp(evt));
        this.resetButton.addEventListener('click', (evt) => this.#handleReset(evt));

        this.eventEmitter.on('ON_RETRIEVE_POSTCODE_LOCATION', () => this.#hideErrorMessage());
    }

    /**
     * Handles the form reset.
     * @param evt
     */
    #handleReset(evt) {
        evt.preventDefault();

        this.eventEmitter.fire('ON_USER_RESET_MAP');
        this.postcodeTextBox.value = '';
        this.resetButton.disabled = 'disabled';
        this.submitButton.disabled = 'disabled'
        this.errorMessage.innerText = ''
        this.errorMessage.style.display = 'none';
    }

    /**
     * Handles the postcode key up event
     * @param evt
     */
    #handlePostCodeKeyUp(evt) {
        const val = evt.target.value;

        this.submitButton.disabled = val ? '': 'disabled';
        this.resetButton.disabled = val ? '': 'disabled';
    }

    /**
     * Form submit handler
     * @param evt - form submit event
     */
    #submitHandler(evt) {
        evt.preventDefault();

        const postcode = evt.target['postalcode'].value;

        this.#fetchPostcode(postcode)
            .then((data) => this.eventEmitter.fire('ON_RETRIEVE_POSTCODE_LOCATION', data))
            .catch((err) => this.#showErrorMessage(err));
    }

    /**
     * Hides the error message
     */
    #hideErrorMessage() {
        this.errorMessage.style.display = 'none';
    }

    /**
     * Handles an error when fetching the postcode.
     * We simply show the error message
     *
     * @param error
     */
    #showErrorMessage(error) {
        this.errorMessage.innerText = error.message;
        this.errorMessage.style.display = 'block';
    }

    /**
     * Gets the postcode from the backend
     * @param postcode
     * @returns {Promise<*>}
     */
    #fetchPostcode = (postcode) => {
        const url = `/geocode?postcode=${postcode}`;

        return fetch(url)
            .then((response) => {
                switch (response.status) {
                    case 200:
                        return response.json();
                    case 404:
                        throw new Error('Postcode not found');
                    case 500:
                        throw new Error('Could not retrieve postcode');
                }
            });
    }
}
