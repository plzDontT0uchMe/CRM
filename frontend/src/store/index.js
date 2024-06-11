import { createStore } from 'vuex'

export default createStore({
    state: {
        profilePhoto: '',
        selectedTrainerProfile: null,
        selectedTab: 'account',
        theme: 'light'
    },
    mutations: {
        updateProfilePhoto(state, newPhoto) {
            state.profilePhoto = newPhoto
        },
        setSelectedTrainerProfile(state, trainer) {
            state.selectedTrainerProfile = trainer
        },
        setSelectedTab(state, tab) {
            state.selectedTab = tab
        },

        toggleTheme(state) {
            state.theme = state.theme === 'light' ? 'dark' : 'light'
        }
    },
    actions: {
        updateProfilePhoto({ commit }, photoUrl) {
            commit('setProfilePhoto', photoUrl)
        },
        resetProfilePhoto({ commit }) {
            commit('RESET_PROFILE_PHOTO')
        }
    },
    getters: {
        profilePhoto: (state) => state.profilePhoto,
        selectedTrainerProfile: (state) => state.selectedTrainerProfile,
        selectedTab: (state) => state.selectedTab
    }
})
