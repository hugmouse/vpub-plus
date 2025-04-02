"use strict";

/**
 * @file Manages date and time formatting across the forum based on user preferences, stores them in localStorage.
 * Also adds preference controls to /account and /update-account pages.
 */

window.addEventListener("DOMContentLoaded", () => {
    const LOCAL_STORAGE_DATE_KEY = "userPrefDateStyle";
    const LOCAL_STORAGE_TIME_KEY = "userPrefTimeStyle";
    const PREVIEW_ELEMENT_ID = "preview-time";
    const ACCOUNT_PAGE_PATHS = ["/account", "/update-account"];
    const DATE_TIME_STYLE_OPTIONS = ["short", "medium", "long", "full"];
    const DEFAULT_STYLE = "medium";

    const timeElements = document.getElementsByTagName("time");
    let currentPreferences = getSavedPreferences();

    /**
     * Retrieves date and time style preferences from localStorage.
     * @returns {{dateStyle: string, timeStyle: string}} The saved preferences or defaults.
     */
    function getSavedPreferences() {
        return {
            dateStyle: localStorage.getItem(LOCAL_STORAGE_DATE_KEY) || DEFAULT_STYLE,
            timeStyle: localStorage.getItem(LOCAL_STORAGE_TIME_KEY) || DEFAULT_STYLE
        };
    }

    /**
     * Saves date style preference to localStorage.
     * @param {string} dateStyle - The selected date style.
     */
    function saveDateStylePreference(dateStyle) {
        localStorage.setItem(LOCAL_STORAGE_DATE_KEY, dateStyle);
        currentPreferences.dateStyle = dateStyle;
    }

     /**
     * Saves time style preference to localStorage.
     * @param {string} timeStyle - The selected date style.
     */
    function saveTimeStylePreference(timeStyle) {
        localStorage.setItem(LOCAL_STORAGE_TIME_KEY, timeStyle);
        currentPreferences.timeStyle = timeStyle;
    }

    /**
     * Formats all <time> elements on the page with the current preferences.
     */
    function formatAllDates() {
        if (!timeElements || timeElements.length === 0) {
            console.log("[server-time-to-local.js] No <time> elements found to format.");
            return;
        }

        try {
            const formatter = new Intl.DateTimeFormat(navigator.language, {
                dateStyle: currentPreferences.dateStyle,
                timeStyle: currentPreferences.timeStyle
            });

            for (const timeElement of timeElements) {
                const dateTimeString = timeElement.getAttribute("datetime");
                if (dateTimeString) {
                    try {
                        const date = new Date(dateTimeString);
                        timeElement.innerText = formatter.format(date);
                    } catch (error) {
                        console.error(`[server-time-to-local.js] Error parsing date string "${dateTimeString}":`, error, timeElement);
                    }
                }
            }
        } catch (error) {
            console.error("[server-time-to-local.js] Error creating Intl.DateTimeFormat or formatting dates:", error, currentPreferences);
        }
    }

    /**
     * Updates the preview element with the current date/time formatted using current preferences.
     */
    function updatePreview() {
        const previewElement = document.getElementById(PREVIEW_ELEMENT_ID);
        if (!previewElement) {
            return;
        }

        try {
            const formatter = new Intl.DateTimeFormat(navigator.language, {
                dateStyle: currentPreferences.dateStyle,
                timeStyle: currentPreferences.timeStyle
            });
            previewElement.innerText = formatter.format(new Date());
        } catch (error) {
            console.error("[server-time-to-local.js] Error updating preview:", error, currentPreferences);
        }
    }

    /**
     * Creates a labeled dropdown select element for date/time style options.
     * @param {string} id - The ID for the select element.
     * @param {string} labelText - The text for the associated label.
     * @param {string} selectedValue - The currently selected value.
     * @returns {{wrapper: HTMLDivElement, selectElement: HTMLSelectElement}} The container div and the select element.
     */
    function createStyleDropdown(id, labelText, selectedValue) {
        const wrapper = document.createElement("div");
        const label = document.createElement("label");
        label.setAttribute("for", id);
        label.textContent = labelText;

        const selectElement = document.createElement("select");
        selectElement.setAttribute("id", id);
        selectElement.setAttribute("name", id);

        DATE_TIME_STYLE_OPTIONS.forEach(optionValue => {
            const option = document.createElement("option");
            option.value = optionValue;
            // Capitalize first letter for display
            option.textContent = optionValue.charAt(0).toUpperCase() + optionValue.slice(1);
            if (optionValue === selectedValue) {
                option.selected = true;
            }
            selectElement.appendChild(option);
        });

        wrapper.appendChild(label);
        wrapper.appendChild(selectElement);
        return {
            wrapper,
            selectElement
        };
    }

    /**
     * Sets up the date/time preference controls on account pages.
     */
    function setupAccountPageControls() {
        const form = document.querySelector("form");
        if (!form) {
            console.error("[server-time-to-local.js] Could not find form on account page to add settings.");
            return;
        }

        const submitButton = form.querySelector('input[type="submit"]');
        if (!submitButton) {
            console.error("[server-time-to-local.js] Could not find submit button on account page to add settings.");
            return;
        }

        const timeSettingsField = document.createElement("div");
        timeSettingsField.className = "field";
        Object.assign(timeSettingsField.style, {
            display: 'flex',
            gap: '1rem',
            flexWrap: 'wrap',
            marginBottom: '1rem'
        });

        const {
            wrapper: dateWrapper,
        } = createStyleDropdown(
            "date-style-pref",
            "Preferred Date Format:",
            currentPreferences.dateStyle
        );
        const {
            wrapper: timeWrapper,
        } = createStyleDropdown(
            "time-style-pref",
            "Preferred Time Format:",
            currentPreferences.timeStyle
        );

        timeSettingsField.appendChild(dateWrapper);
        timeSettingsField.appendChild(timeWrapper);
        form.insertBefore(timeSettingsField, submitButton);

        updatePreview();
    }

    if (ACCOUNT_PAGE_PATHS.includes(window.location.pathname)) {
        setupAccountPageControls();

        document.addEventListener("change", (e) => {
            const target = e.target;
            let changed = false;
            
            if (target.id === 'date-style-pref') {
                saveDateStylePreference(target.value);
                changed = true;
            } else if (target.id === 'time-style-pref') {
                saveTimeStylePreference(target.value);
                changed = true;
            }
        
            if (changed) {
                updatePreview();
            }
        });
    }

    formatAllDates();
});
